package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/kajikentaro/spotify-file-based-client/client/models"
	"github.com/zmb3/spotify/v2"
)

func FetchRemotePlaylistContent(client *spotify.Client, ctx context.Context) ([]models.PlaylistContent, error) {
	result := []models.PlaylistContent{}
	playlists, err := client.CurrentUsersPlaylists(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range playlists.Playlists {
		content := models.SimplePlaylistToContent(v)
		result = append(result, content)
	}

	return result, nil
}

func FetchRemotePlaylistTrack(client *spotify.Client, ctx context.Context, id string) ([]models.TrackContent, error) {
	LIMIT := 100
	result := []models.TrackContent{}
	for offset := 0; true; offset += LIMIT {
		playlistItemPage, err := client.GetPlaylistItems(ctx, spotify.ID(id), spotify.Limit(LIMIT), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch playlist %s: %s", id, err)
		}
		for _, playlistItem := range playlistItemPage.Items {
			track := playlistItem.Track.Track
			trackContent := models.FullTrackToContent(track)
			result = append(result, trackContent)
		}
		if len(playlistItemPage.Items) != LIMIT {
			break
		}
	}
	return result, nil
}

func CreateRemotePlaylist(client *spotify.Client, ctx context.Context, name string) (models.PlaylistContent, error) {
	user, err := client.CurrentUser(ctx)
	if err != nil {
		return models.PlaylistContent{}, fmt.Errorf("failed to get a current user info: %w", err)
	}
	new, err := client.CreatePlaylistForUser(ctx, user.ID, name, "", false, false)
	if err != nil {
		return models.PlaylistContent{}, fmt.Errorf("failed to create playlist %s: %w", name, err)
	}
	return models.SimplePlaylistToContent(new.SimplePlaylist), nil
}

type EditTrackRes struct {
	models.TrackContent
	IsOk    bool
	Message string
}

func AddRemoteTrack(client *spotify.Client, ctx context.Context, playlistId string, tracks []models.TrackContent) ([]EditTrackRes, error) {
	if playlistId == "" {
		return []EditTrackRes{}, fmt.Errorf("playlistId is empty")
	}

	result := []EditTrackRes{}
	trackIds := []spotify.ID{}
	for _, v := range tracks {
		if v.Id != "" {
			result = append(result, EditTrackRes{v, true, ""})
			trackIds = append(trackIds, spotify.ID(v.Id))
		} else {
			// IDがないときは検索する
			res, err := client.Search(ctx, v.SearchQuery(), spotify.SearchTypeTrack, spotify.Limit(1))
			if err != nil {
				return []EditTrackRes{}, fmt.Errorf("failed to search: %w", err)
			}
			if len(res.Tracks.Tracks) == 0 {
				result = append(result, EditTrackRes{v, false, "no search results found"})
				continue
			}
			trackIds = append(trackIds, res.Tracks.Tracks[0].ID)
			content := models.FullTrackToContent(&res.Tracks.Tracks[0])
			content.FileName = v.FileName
			result = append(result, EditTrackRes{content, true, ""})
		}
	}
	if len(trackIds) == 0 {
		// 検索結果が何も見つからなかった場合
		return result, nil
	}

	LIMIT := 100
	for offset := 0; true; offset += LIMIT {
		if len(trackIds)-1 < offset {
			break
		}
		var trackChunk []spotify.ID
		if offset+LIMIT < len(trackIds) {
			trackChunk = trackIds[offset : offset+LIMIT]
		} else {
			trackChunk = trackIds[offset:]
		}
		_, err := client.AddTracksToPlaylist(ctx, spotify.ID(playlistId), trackChunk...)
		if err != nil {
			return []EditTrackRes{}, err
		}
	}
	return result, nil
}

func RemoveRemoteTrack(client *spotify.Client, ctx context.Context, playlist models.PlaylistContent, tracks []models.TrackContent) error {
	if playlist.Id == "" {
		return errors.New("playlist id is empty")
	}

	trackIds := []spotify.ID{}
	for _, v := range tracks {
		if v.Id == "" {
			return fmt.Errorf("track %v is not have track id", v)
		}
		trackIds = append(trackIds, spotify.ID(v.Id))
	}

	_, err := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlist.Id), trackIds...)
	if err != nil {
		return err
	}

	return nil
}

func RemoveRemotePlaylist(client *spotify.Client, ctx context.Context, playlist models.PlaylistContent) error {
	if playlist.Id == "" {
		return errors.New("playlist id is empty")
	}
	err := client.UnfollowPlaylist(ctx, spotify.ID(playlist.Id))
	if err != nil {
		return err
	}
	return nil
}
