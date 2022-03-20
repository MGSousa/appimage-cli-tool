package install

import (
	"appimage-cli-tool/cmd/commands"
	"appimage-cli-tool/internal/repos"
	"appimage-cli-tool/internal/utils"

	"fmt"
	"os"
)

type InstallCmd struct {
	Target string `arg name:"target" help:"Installation target." type:"string"`
}

func (cmd *InstallCmd) Run(ctx *commands.Context) (err error) {
	if _, err := os.Stat(cmd.Target); err == nil {
		cmd.createDesktopIntegration(cmd.Target)
		return nil
	}

	repo, err := repos.ParseTarget(cmd.Target)
	if err != nil {
		return err
	}

	release, err := repo.GetLatestRelease()
	if err != nil {
		return err
	}

	selectedBinary, err := utils.PromptBinarySelection(release.Files)
	if err != nil {
		return err
	}

	targetFilePath, err := utils.MakeTargetFilePath(selectedBinary)
	if err != nil {
		return err
	}

	if _, err = os.Stat(targetFilePath); err == nil {
		return ApplicationInstalled
	}

	err = repo.Download(selectedBinary, targetFilePath)
	if err != nil {
		return err
	}

	cmd.addToRegistry(targetFilePath, repo)

	if ctx.DesktopIntegration {
		cmd.createDesktopIntegration(targetFilePath)
	}

	signingEntity, _ := utils.VerifySignature(targetFilePath)
	if signingEntity != nil {
		fmt.Println("AppImage signed by:")
		for _, v := range signingEntity.Identities {
			fmt.Println("\t", v.Name)
		}
	}

	return
}

func (cmd *InstallCmd) addToRegistry(targetFilePath string, repo repos.Repo) {
	sha1, _ := utils.GetFileSHA1(targetFilePath)
	updateInfo, _ := utils.ReadUpdateInfo(targetFilePath)
	if updateInfo == "" {
		updateInfo = repo.FallBackUpdateInfo()
	}

	entry := utils.RegistryEntry{
		FilePath:   targetFilePath,
		Repo:       repo.Id(),
		FileSha1:   sha1,
		UpdateInfo: updateInfo,
	}

	registry, _ := utils.OpenRegistry()
	if registry != nil {
		_ = registry.Add(entry)
		_ = registry.Close()
	}
}

func (cmd *InstallCmd) createDesktopIntegration(targetFilePath string) {
	libAppImage, err := utils.NewLibAppImageBindings()
	if err != nil {
		fmt.Println("Integration failed: missing libappimage.so")
		return
	}

	fmt.Println("Creating menu entry and mime-types integrations")
	err = libAppImage.Register(targetFilePath)
	if err != nil {
		fmt.Println("Integration failed: " + err.Error())
	} else {
		fmt.Println("Integration completed")
	}
}
