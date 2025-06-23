let time = 0;
let screen_width = 128;
let screen_height = 128;
let A_speed = 0.05;
let B_speed = 0.03;
let R1 = 1.0;
let R2 = 2.0;
let K1 = 5.0;
let K2 = 5.0;
let cx = screen_width / 2;
let cy = screen_height / 2;

let _draw = fn() {
    for i = 0, 127 {
        for j = 0, 127 {
            pset(i, j, 233, 0, 0);
        }
    }
}

let _update = fn() {
    pset(0, 0, 33, 33, 33);
}