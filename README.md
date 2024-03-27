# ACOUSTIC
Listen, Control and See the details of your favorite songs.

# How To Use
`$ go run . your_directory/`

![alt text](images/Screenshot%20from%202024-03-27%2002-37-23.png)

# Options
You can change the default characters that are used to print the song image a progress bar. Just specify imageChar and progressbarChar flags in your command:

`$ go run . -image-char='your_char' -progressbar-char='your_char' your_directory/`

![alt text](images/Screenshot%20from%202024-03-27%2002-44-06.png)

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
| q           | Exit the program          |

# Roadmap
- Fix possible memory leak
- Fix seek bug that does not update the progress bar
- Create terminal command
- Replace deprecated packages
- Add keyboard key settings
- Add Cool/Cold feature
- Add Auto-Pause functionality