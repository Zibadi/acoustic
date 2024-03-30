# ACOUSTIC
Listen, Control and See the details of your favorite musics.

# How To Run
`$ go run . your_directory/`

![alt text](images/Screenshot%20from%202024-03-27%2002-37-23.png)

# Settings
You can change the default characters that are used to print the music image and progress bar. Just specify imageChar and progressbarChar flags in your command:

`$ go run . -imageChar='your_char' -progressbarChar='your_char' your_directory/`

![alt text](images/Screenshot%20from%202024-03-27%2002-44-06.png)

# Auto-Pause
If you open a video or another audio file on your computer while the music is playing, the program will automatically pause the music. After the external audio gets finished, the music will continue to play. So cool, Right!!!

# Playback Controls
During playing the music, you can control the player. Press these keyboard keys to control the player:

| Key         | Action                    |
| ----------- | ------------------------- |
| n           | Go to next music          |
| p           | Go to previous music      |
| space       | Pause/Play                |
| UpArrow     | Increase the sound volume |
| DownArrow   | Decrease the sound volume |
| RightArrow  | Seek forward              |
| LeftArrow   | Seek backward             |
| s           | Shuffle the play list     |
| q           | Exit the program          |

# Roadmap
- Create terminal command
- Add keyboard key settings