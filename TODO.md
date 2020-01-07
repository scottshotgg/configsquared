# TODO

1.  pull out folder names, rework assets
2.  file names as flags at compile time of cc
3.  Inject the assets into the binary
4.  required and default should be used together
5.  fix/re/move some of the panics
6.  Make "Must\*" functions
7.  Rework template file
8.  Add tests
9.  Make a better parser
10. Change template.go to be based on files
11. Switch to actual code generation with github.com/moznion/gowrtr or similar
12. Handle nested struct variables
13. Bootstrap arguments using this lib
14. Figure out how to make bools not NEED an argument
15. Implement outstanding types:

- `array`
- `slice`
- `object`
- ~~`date`~~ - just use `time` if you want to do this
- ~~`byte`~~
- ~~`unix time`~~
- ~~`duration`~~
- ~~`any`~~
- ~~`time`~~
- ~~`float32/64`~~
- ~~`int32/64`~~
- ~~`uint32/64`~~
- ~~`ip`~~
- ~~`ipv4`~~
- ~~`ipv6`~~
- ~~`url`~~

16. Could we allow the user to implement their own types?
    - make it very pluggable
17. Allow VERY rudimentary type specifications:

- lt
- gt
- eq
- ne
- max
- min
- etc

18. ~~Add `example` as a valid key~~
19. How to make our own help text
