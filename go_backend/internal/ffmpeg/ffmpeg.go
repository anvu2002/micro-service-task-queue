package ffmpeg

import (
	"fmt"
	"os/exec"
)

func CreateFrame(
	save_path string, audio_path string, image_path string, vfilter string) []byte {
	if vfilter == "" {
		vfilter = "pad=ceil(iw/2)*2:ceil(ih/2)*2"
	}
	command := fmt.Sprintf("ffmpeg -i %s -i %s -vf \"%s\" %s ", audio_path, image_path, vfilter, save_path)

	out, err := exec.Command(command).Output()
	fmt.Print(err)
	return out
}

func Concat(save_path string, frames_path string) []byte {
	command := fmt.Sprintf("ffmpeg -f concat -safe 0 -i %s %s", frames_path, save_path)
	out, err := exec.Command(command).Output()
	fmt.Print(err)
	return out
}

func Test(audio_path string, images string) (string, error) {
	res := string("FFMPEG endpoint is working!")

	return res, nil
}
