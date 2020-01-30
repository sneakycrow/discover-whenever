package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zmb3/spotify"
	"io/ioutil"
	"os"
)

func main() {
	// Initialize flags
	var playlistName string
	var playlistDesc string
	flag.StringVar(&playlistName, "name", "Manually Generated Playlist by Sneaky Crow", "Playlist name")
	flag.StringVar(&playlistDesc, "desc", "A randomly generated playlist by sneakycrow's playlist generator", "Playlist description")
	flag.Parse()
	// Login and grab the client
	client := login()
	// Grab the user
	user, err := client.CurrentUser()
	// Create a time range typed string, we'll use this later
	timeRange := "short"
	// Grab our top 5 artists based on short term time range
	topArtistOps := spotify.Options{
		Timerange: &timeRange,
	}
	topArtists, err := client.CurrentUsersTopArtistsOpt(&topArtistOps)
	check(err)
	var TopArtistIDs [5]spotify.ID
	for index, artist := range topArtists.Artists[:5] {
		TopArtistIDs[index] = artist.ID
	}
	// grab our top 5 tracks based on short term time range
	topTracksOps := spotify.Options{
		Timerange: &timeRange,
	}
	topTracks, err := client.CurrentUsersTopTracksOpt(&topTracksOps)
	check(err)
	var TopTracksIDs [5]spotify.ID
	for index, track := range topTracks.Tracks[:5] {
		TopTracksIDs[index] = track.ID
	}
	// Create a Seeds type of 2 tracks and 3 artists (maximum seeds we can provide is 5)
	reccSeeds := spotify.Seeds{
		Artists: TopArtistIDs[:2],
		Tracks:  TopTracksIDs[:3],
	}
	// Grab our recommendations
	recommendations, err := client.GetRecommendations(reccSeeds, nil, nil)
	check(err)
	// Create an array of type spotify.ID of TrackIDs
	var RecommendedTracks []spotify.ID
	for _, recc := range recommendations.Tracks {
		RecommendedTracks = append(RecommendedTracks, recc.ID)
	}
	// Create a new playlist
	newPlaylist, err := client.CreatePlaylistForUser(user.ID, playlistName, playlistDesc, true)
	check(err)
	_, err = client.AddTracksToPlaylist(newPlaylist.ID, RecommendedTracks...)
	check(err)
	fmt.Printf("New Playlist Created: %s!", playlistName)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func outputTopArtistsAsJSON(artists []spotify.FullArtist) {
	marshaledTopArtists, err := json.Marshal(artists)
	check(err)
	cwdPath, err := os.Getwd()
	check(err)
	err = ioutil.WriteFile(fmt.Sprintf("%s/output.json", cwdPath), marshaledTopArtists, 0644)
	check(err)
}
