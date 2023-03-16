package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	files := []File{
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
		{path: "https://www.google.com/url?sa=i&url=https%3A%2F%2Fpicsum.photos%2Fimages&psig=AOvVaw3khh4SV9yAjz6VCZ26iQGh&ust=1678878455630000&source=images&cd=vfe&ved=0CA8QjRxqFwoTCPCGhuGj2_0CFQAAAAAdAAAAABAE"},
	}

	readers := downloadParallel(files)

	name := []string{"dwnload.zip"}

	err := Archive(name, readers...)
	if err != nil {
		panic(err)
	}

}

func downloadParallel(files []File) []io.Reader {
	start := time.Now()

	ch := make(chan io.Reader)
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file File) {
			defer wg.Done()
			read, err := file.Download()
			if err != nil {
				panic(err)
			}
			ch <- read
		}(file)
	}

	reader := []io.Reader{}
	go func() {
		for r := range ch {
			reader = append(reader, r)
		}
	}()

	wg.Wait()
	fmt.Println("time elapsed ", time.Since(start), " for readers ", len(reader))
	return reader

}

func downloadSerial(files []File) []io.Reader {
	start := time.Now()

	reader := []io.Reader{}
	for _, file := range files {
		read, err := file.Download()
		if err != nil {
			panic(err)
		}
		reader = append(reader, read)
	}

	fmt.Println("time elapsed ", time.Since(start), " for readers ", len(reader))
	return reader

}

type File struct {
	path string
}
type Downloader interface {
	download(url string) (r io.Reader, err error)
}

type Archiver interface {
	Archive(names []string, readers ...io.Reader) error
}

type zipp struct {
}

func (f *File) Download() (io.Reader, error) {
	fmt.Println(fmt.Sprintf("%s being downloaded", f.path))
	resp, err := http.Get(f.path)
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(out)
	return reader, nil
}

func Archive(names []string, readers ...io.Reader) error {
	out, err := os.Create(names[0])
	if err != nil {
		return err

	}
	zipW := zip.NewWriter(out)
	defer out.Close()

	for i := 0; i < len(readers); i++ {

		w1, err := zipW.Create(fmt.Sprintf("%d", rand.Intn(10)))
		if err != nil {
			return err
		}
		reader := readers[i]
		_, err = io.Copy(w1, reader)
		if err != nil {
			return err
		}

		defer zipW.Close()
	}
	return nil

}

// package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// )

// type Downloader interface {
// 	Download(url string) (r io.Reader, err error)
// }

// type Archiver interface {
// 	Archive(names []string, readers ...io.Reader) (outR io.Reader, err error)
// }

// type zip struct{}

// type Video struct {
// 	Id int
// }

// type web struct {
// }

// func main() {
// 	//fileUrl := "https://filesamples.com/samples/video/mp4/sample_1280x720_surfing_with_audio.mp4"

// 	list := []string{}
// 	list = append(list, "https://filesamples.com/samples/video/mp4/sample_1280x720_surfing_with_audio.mp4")

// 	list = append(list, "https://filesamples.com/samples/video/mp4/sample_960x400_ocean_with_audio.mp4")
// 	//err := fileDownload("assignment3_video.mp4", fileUrl)

// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	//downloader:=web.NewDownloader()

// 	var video Video

// 	responseList := []io.Reader{}

// 	for i, x := range list {
// 		//fmt.Println(i)
// 		resp, err := video.Download(x, fmt.Sprintf("video%d.mp4", i))

// 		responseList = append(responseList, resp)
// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Println(fmt.Sprintf("Downloaded video is video%d", i))

// 	}

// 	var zipper zip

// 	zipR,err:=zipper.Archive()

// }

// // func Newfile(path string) string {
// // 	return path + ".mp4"
// // }

// func (v Video) Download(url, filepath string) (r io.Reader, err error) {
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	//buf := make([]byte, 1024)

// 	//filepath := Newfile(fmt.Sprintf("video%d", v.Id))

// 	output, err := os.Create(filepath)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer output.Close()

// 	_, err = io.Copy(output, response.Body)

// 	return response.Body, err

// }

// func (z zip) Archive(names []string, readers ...io.Reader) (outR []io.Reader, err error){

// }

// // func fileDownload(filepath, url string) error {
// // 	response, err := http.Get(url)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	defer response.Body.Close()

// // 	output, err := os.Create(filepath)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer output.Close()

// // 	_, err = io.Copy(output, response.Body)

// // 	return err

// // }
