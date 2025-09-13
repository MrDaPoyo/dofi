WIDTH = 128
HEIGHT = 128

function init()
end

function update()
end

function draw()
    for x = 0, WIDTH - 1 do
        for y = 0, HEIGHT - 1 do
            local value = ((x + y) % 2) * 255
            gfx.pset(x, y, value, value, value)
        end
    end
end
