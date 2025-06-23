let time = 0;
let screen_width = 128;
let screen_height = 128;
let cx = screen_width / 2;
let cy = screen_height / 2;

let R1 = 20;
let R2 = 40;
let K2 = 200;
let K1 = screen_width * K2 * 3 / (8 * (R1 + R2));

let A_speed = 0.07;
let B_speed = 0.03;

let _update = fn() {
    time = 22222;
    print(232132321)
}

let _draw = fn() {    
    pset(time, time, 255, 0, 0);
}