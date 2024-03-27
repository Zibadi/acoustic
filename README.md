# ACOUSTIC
Listen, Control and See the details of your favorite songs.

# How To Run
`$ go run . your_directory/`

![alt text](images/Screenshot%20from%202024-03-27%2002-37-23.png)

# Settings
You can change the default characters that are used to print the song image a progress bar. Just specify imageChar and progressbarChar flags in your command:

`$ go run . -imageChar='your_char' -progressbarChar='your_char' your_directory/`

![alt text](images/Screenshot%20from%202024-03-27%2002-44-06.png)

# Auto-Puase
If you open a video or another audio file on your computer while the music is playing, the program will automatically pause the music. After the external audio gets finished, the music will continue to play. So cool, Right!!!

# Playback Controls
During playing the song, you can control the player. Press these keyboard keys to control the player:

| Key         | Action                    |
| ----------- | ------------------------- |
| n           | Go to next song           |
| p           | Go to previous song       |
| space       | Pause/Play                |
| UpArrow     | Increase the sound volume |
| DownArrow   | Decrease the sound volume |
| RightArrow  | Seek forward              |
| LeftArrow   | Seek backward             |
| s           | Shuffle the play list     |
| q           | Exit the program          |

# Roadmap
- Fix seek bug that does not update the progress bar
- Create terminal command
- Replace deprecated packages
- Add keyboard key settings
- Add Cool/Cold feature