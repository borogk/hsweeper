![](logo.png)

_H-Sweeper_ is a Minesweeper clone with extra lives mechanic, that runs entirely in terminal.
Extra lives were added because it's irritating to play well, only to encounter a pure 50-50 that ruins your run.

### How to install

> [!IMPORTANT]
> Requires [Go version 1.25 or newer](https://go.dev/doc/install)

Install by running in terminal:
```shell
git clone https://github.com/borogk/hsweeper.git
cd hsweeper
go install
```

After installing, simply run:
```shell
hsweeper
```

### How to play

You are expected to already know how to play Minesweeper.

In case you don't know, [Minesweeper article on Wikipedia](https://en.wikipedia.org/wiki/Minesweeper_(video_game))
does a decent job explaining. _H-Sweeper_ doesn't change rules besides adding extra lives, which is explained below.

### Game modes

Default game mode is _H-Expert_, which plays the same as regular Minesweeper Expert mode, but with one extra life.

_H-Big_ extends _H-Expert_ by having a bigger playing field and giving ever more extra lives to compensate.

_Classic_ modes play exactly like the 3 modes of Windows Minesweeper with no extra lives.

| # | Mode           | Size                                           | Mines                         | Extra lives        |
|---|----------------|------------------------------------------------|-------------------------------|--------------------|
| 1 | H-Expert       | 30 x 16                                        | 99                            | +1                 |
| 2 | H-Big          | Fits entire terminal window (at least 30 x 16) | 1 every 5 cells (at least 99) | +1 every 480 cells |
| 3 | Classic Easy   | 9 x 9                                          | 10                            | None               |
| 4 | Classic Medium | 16 x 16                                        | 40                            | None               |
| 5 | Classic Expert | 30 x 16                                        | 99                            | None               |

### ♥ Extra lives ♥

Current amount of lives is represented by `♥` symbols in the top left corner.

If you have more than one, revealing a bomb loses one life instead of losing the game.
Exploded bomb is removed from the game and adjacent numbers are adjusted accordingly.

Extra lives are not given immediately, but are rather rewarded for revealing some amount of play field.
Only on huge game sizes (2400 cells and above) a few are granted right away.

> [!TIP]
> Extra lives are supposed to help in absolute uncertainty! Try solving as much as you can without relying on them.

### Controls

The game is controlled only with keyboard, there is no mouse support.

| Key                 | Function                                               |
|---------------------|--------------------------------------------------------|
| `1` `2` `3` `4` `5` | Select title menu option                               |
| `←` `→` `↑` `↓`     | Move the cursor                                        |
| `SPACE` `↵`         | _Action Key_ (explained below)                         |
| `R`                 | Reveal cell (must not be marked)                       |
| `F`                 | Toggle flag mark `⚑`                                   |
| `Q`                 | Toggle question mark `?`                               |
| `DELETE` `⌫`        | Clear `⚑` or `?`                                       |
| `ESC`               | Quits to title menu (if already there, quits the game) |
| `CTRL-C`            | Quits the game                                         |

### Action Key

_Action Key_ is context-sensitive, combining functions of
`Left-click`, `Right-click` and `Left+right-click` of Windows Minesweeper.

| Condition                                                 | Function                                             |
|-----------------------------------------------------------|------------------------------------------------------|
| In title menu                                             | Start new H-Expert game                              |
| First move of a game                                      | Reveal around the cursor and start for real          |
| On unrevealed cell                                        | Toggle `⚑`                                           |
| On cells, where amount of adjacent `⚑` matches the number | Reveal unmarked adjacent cells (known as "chording") |
| On cells with `♥`                                         | Pick up extra life                                   |
| After game over                                           | Restart game                                         |

> [!NOTE]
> Design ideas behind having this control scheme:
> 1. Mostly play with a single button (other than moving the cursor of course).
> 2. Minimize risk of losing by accident, as it puts flags on unrevealed cells, rather than revealing them.
> 3. Force-revealing cells requires a more conscious decision to press separate `R` button.

### Future support notice

Future support and updates are unlikely. This was originally just my personal project to learn Go programming language.
The game is minimalistic, made with personal tastes in mind and with little to no customization options.

### Final note

> [!NOTE]
> H stands for H

### Author

Originally created by **borogk** in 2025.
