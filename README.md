# Youtube Channel Downloader

Simple CLI tool written in Go to download all videos of a YouTube channel up to 500 videos. It makes use of YouTube Data API to list all the videos of a channel and downloads them using package [kkdai/youtube](https://github.com/kkdai/youtube).

## Configuration
Two environment variables can be used to configure, one of them is necessary. 

* YT_API_KEY (Youtube API key)
* YT_OUTPUT_DIR (Output directory, optional, defaults to user's home directory)

## Usage

Simply pass YouTube channel ID (which is the last part of the channel URL) as channel flag.
As an example, for channel https://www.youtube.com/channel/UC68ldbHwL_-5qzETqOaAMWQ `UC68ldbHwL_-5qzETqOaAMWQ` is the ID so that would be:
```
./yt-channel-downloader -channel UC68ldbHwL_-5qzETqOaAMWQ
```

This will download each video of the given channel in 1080p and if does not exist, it will check 720p quality. Quality below that is not currently supported, videos which do not have at least 720p will be skipped.
