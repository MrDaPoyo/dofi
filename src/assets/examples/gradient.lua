WIDTH = 128
HEIGHT = 128

function init()
end

function update()
end

function draw()
    for x = 0, WIDTH - 1 do
        for y = 0, HEIGHT - 1 do
            local r = math.floor((x / (WIDTH - 1)) * 255)
            local g = 0
            local b = 0
            gfx.pset(x, y, r, g, b)
        end
    end
end