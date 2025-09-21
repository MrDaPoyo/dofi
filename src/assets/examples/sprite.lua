local W, H = 128, 128
local t = 0

function init()
	t = 0
end

function update()
	t = t + 1
end

function draw()
	for x = 0, W - 1 do
		for y = 0, H - 1 do
			local shade = (x + y + t) % 64
			gfx.pset(x, y, 0, shade, 16)
		end
	end

	local px = math.floor((math.sin(t * 0.04) * 0.5 + 0.5) * (W - 8))
	local py = math.floor((math.sin(t * 0.07 + 1.2) * 0.5 + 0.5) * (H - 8))

	gfx.sprite(0, px, py)
	gfx.sprite(0, (px + 4) % (W - 8), (py + 2) % (H - 8))
end
