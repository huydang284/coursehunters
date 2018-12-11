// Sample Go code for user authorization

package main

import (
    "encoding/json"
    "fmt"
    "golang.org/x/oauth2/google"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/user"
    "path/filepath"
    "strconv"
    "strings"

    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "google.golang.org/api/youtube/v3"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
    cacheFile, err := tokenCacheFile()
    if err != nil {
        log.Fatalf("Unable to get path to cached credential file. %v", err)
    }
    tok, err := tokenFromFile(cacheFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(cacheFile, tok)
    }
    return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
        "authorization code: \n%v\n\n", authURL)
    fmt.Print(">> Enter your code: ")

    var code string
    if _, err := fmt.Scan(&code); err != nil {
        log.Fatalf("Unable to read authorization code %v", err)
    }

    tok, err := config.Exchange(oauth2.NoContext, code)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web %v", err)
    }
    return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
    usr, err := user.Current()
    if err != nil {
        return "", err
    }
    tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
    os.MkdirAll(tokenCacheDir, 0700)
    return filepath.Join(tokenCacheDir,
        url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    t := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(t)
    defer f.Close()
    return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", file)
    f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    errEncode := json.NewEncoder(f).Encode(token)
    if errEncode != nil {
        log.Fatal(errEncode)
    }
}

func handleError(err error, message string) {
    if message == "" {
        message = "Error making API call"
    }
    if err != nil {
        log.Fatalf(message+": %v", err.Error())
    }
}

func createPlaylist(service *youtube.Service, playlistName string) string {
    log.Println("Creating playlist with title: " + playlistName)
    properties := map[string]string{"snippet.title": playlistName}
    res := createResource(properties)

    return playlistsInsert(service, "snippet", res)
}

func uploadVideo(service *youtube.Service, filePath string, lessonTitle string) string {
    log.Println("Uploading " + filePath)
    upload := &youtube.Video{
        Snippet: &youtube.VideoSnippet{
            Title:      lessonTitle,
            CategoryId: "27",
        },
        Status: &youtube.VideoStatus{PrivacyStatus: "private"},
    }
    call := service.Videos.Insert("snippet,status", upload)
    file, err := os.Open(filePath)
    defer file.Close()
    if err != nil {
        log.Fatalf("Error opening %v: %v", "sample.mp4", err)
    }

    response, err := call.Media(file).Do()
    if err != nil {
        log.Fatalf("Error making YouTube API call: %v", err)
    }
    log.Printf("Uploaded successful! Video ID: %v\n", response.Id)

    return response.Id
}

func createResource(properties map[string]string) string {
    resource := make(map[string]interface{})
    for key, value := range properties {
        keys := strings.Split(key, ".")
        ref := addPropertyToResource(resource, keys, value, 0)
        resource = ref
    }
    propJson, err := json.Marshal(resource)
    if err != nil {
        log.Fatal("cannot encode to JSON ", err)
    }
    return string(propJson)
}

func addPropertyToResource(ref map[string]interface{}, keys []string, value string, count int) map[string]interface{} {
    for k := count; k < (len(keys) - 1); k++ {
        switch val := ref[keys[k]].(type) {
        case map[string]interface{}:
            ref[keys[k]] = addPropertyToResource(val, keys, value, k+1)
        case nil:
            next := make(map[string]interface{})
            ref[keys[k]] = addPropertyToResource(next, keys, value, k+1)
        }
    }
    // Only include properties that have values.
    if count == len(keys)-1 && value != "" {
        valueKey := keys[len(keys)-1]
        if valueKey[len(valueKey)-2:] == "[]" {
            ref[valueKey[0:len(valueKey)-2]] = strings.Split(value, ",")
        } else if len(valueKey) > 4 && valueKey[len(valueKey)-4:] == "|int" {
            ref[valueKey[0:len(valueKey)-4]], _ = strconv.Atoi(value)
        } else if value == "true" {
            ref[valueKey] = true
        } else if value == "false" {
            ref[valueKey] = false
        } else {
            ref[valueKey] = value
        }
    }
    return ref
}

func playlistsInsert(service *youtube.Service, part string, res string) string {
    resource := &youtube.Playlist{}
    if err := json.NewDecoder(strings.NewReader(res)).Decode(&resource); err != nil {
        log.Fatal(err)
    }
    call := service.Playlists.Insert(part, resource)
    resp, err := call.Do()
    handleError(err, "")

    return resp.Id
}

func addVideoToPlaylist(service *youtube.Service, playlistId string, videoId string) {
    properties := map[string]string{"snippet.playlistId": playlistId,
        "snippet.resourceId.kind":    "youtube#video",
        "snippet.resourceId.videoId": videoId,
    }
    res := createResource(properties)
    playlistItemsInsert(service, "snippet", res)
}

func playlistItemsInsert(service *youtube.Service, part string, res string) {
    resource := &youtube.PlaylistItem{}
    if err := json.NewDecoder(strings.NewReader(res)).Decode(&resource); err != nil {
        log.Fatal(err)
    }
    call := service.PlaylistItems.Insert(part, resource)
    _, err := call.Do()
    handleError(err, "")
}

func initService() *youtube.Service {
    ctx := context.Background()

    b, err := ioutil.ReadFile("client_secret.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
    }

    client := getClient(ctx, config)
    service, err := youtube.New(client)

    handleError(err, "Error creating YouTube client")

    return service
}

func getPlaylist(service *youtube.Service, playlistName string) string {
    // get all playlists
    call := service.Playlists.List("snippet").Mine(true)
    response, err := call.Do()
    handleError(err, "")
    for _, item := range response.Items {
        if item.Snippet.Title == playlistName {
            log.Println("Created playlist found: " + item.Id)
            return item.Id
        }
    }

    log.Println("Cannot find the playlist")
    return createPlaylist(service, playlistName)
}
