package models

import "testing"

func Test_TrackContent_marshal(t *testing.T) {
	tc := TrackContent{
		Id:       "123",
		Name:     "test track",
		Artist:   "test artist",
		Album:    "test album",
		Seconds:  "123",
		Isrc:     "ABCDEFG",
		FileName: "test track.txt",
	}
	actual := tc.marshal()
	expected :=
		`id 123
name test track
artist test artist
album test album
seconds 123
isrc ABCDEFG
file_name test track.txt
`
	if expected != actual {
		t.Errorf("\nactual:\n%s \nexpected:\n%s", actual, expected)
	}
}

func Test_PlaylistContent_marshal(t *testing.T) {
	tc := PlaylistContent{
		Id:      "123",
		Name:    "test playlist name",
		DirName: "test playlist name",
	}
	actual := tc.marshal()
	expected :=
		`NOTE: Do not delete or edit this file.

id 123
name test playlist name
dir_name test playlist name
`
	if expected != actual {
		t.Errorf("\nactual:\n%s \nexpected:\n%s", actual, expected)
	}
}

func Test_unmarshalPlaylistContent(t *testing.T) {
	text :=
		`id 123
name test playlist name
dir_name test playlist name
`

	actual := unmarshalPlaylistContent(text)
	expected := PlaylistContent{
		Id:      "123",
		Name:    "test playlist name",
		DirName: "test playlist name",
	}
	if expected != actual {
		t.Errorf("\nactual:\n%s \nexpected:\n%s", actual, expected)
	}
}

func Test_unmarshalTrackContent(t *testing.T) {
	text :=
		`id 123
name test track
artist test artist
album test album
seconds 123
isrc ABCDEFG
file_name test track.txt
`

	actual := unmarshalTrackContent(text)
	expected := TrackContent{
		Id:       "123",
		Name:     "test track",
		Artist:   "test artist",
		Album:    "test album",
		Seconds:  "123",
		Isrc:     "ABCDEFG",
		FileName: "test track.txt",
	}
	if expected != actual {
		t.Errorf("\nactual:\n%s \nexpected:\n%s", actual, expected)
	}
}