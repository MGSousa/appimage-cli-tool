package repos

import (
	"appimage-cli-tool/internal/utils"

	"encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "net/url"
        "strings"
)

type (
	AppImageHubRepo struct {
		ContentId string
	}
	Store struct {
		Id   string
		Name string
		Url  string
	}
)

func NewAppImageHubRepo(target string) (Repo, error) {
	if strings.HasPrefix(target, "https://www.appimagehub.com/p/") {
		target = strings.Replace(target, "https://www.appimagehub.com/p/", "appimagehub:", 1)
	}

	if !strings.HasPrefix(target, "appimagehub:") {
		return nil, InvalidTargetFormat
	}

	repo := &AppImageHubRepo{}
	repo.ContentId = target[12:]

	return repo, nil
}

func (a AppImageHubRepo) Id() string {
	return "appimagehub:" + a.ContentId
}

func (a AppImageHubRepo) GetLatestRelease() (*Release, error) {
	var (
                downloadLinks []utils.BinaryUrl
                link          string
        )
	store := []Store{}

        req, err := http.Get(fmt.Sprintf("https://www.appimagehub.com/p/%s/loadFiles", a.ContentId))
        if err != nil {
                return nil, err
        }
        defer req.Body.Close()

        content, _ := ioutil.ReadAll(req.Body)
        if err := json.Unmarshal(content, &store); err != nil {
                return nil, err
        }
        if len(store) > 0 {
                for _, v := range store {
                        if link, err = url.QueryUnescape(v.Url); err != nil {
                                return nil, err
                        }
                        downloadLink := utils.BinaryUrl{
                                FileName: v.Name,
                                Url:      link,
                        }
                        if strings.HasSuffix(downloadLink.FileName, ".AppImage") ||
                                strings.HasSuffix(downloadLink.FileName, ".appimage") {
                                downloadLinks = append(downloadLinks, downloadLink)
                        }
                }
        }

	if len(downloadLinks) > 0 {
		return &Release{
			"latest",
			downloadLinks,
		}, nil
	} else {
		return nil, NoAppImageBinariesFound
	}
}

func (a AppImageHubRepo) Download(binaryUrl *utils.BinaryUrl, targetPath string) (err error) {
	err = utils.DownloadAppImage(binaryUrl.Url, targetPath)
	return
}

func (a AppImageHubRepo) FallBackUpdateInfo() string {
	return "ocs-v1-appimagehub-zsync|www.appimagehub.com/ocs/v1|" + a.ContentId + "|*.AppImage"
}
