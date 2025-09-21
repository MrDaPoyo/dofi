local W, H = 128, 128
local t = 0

function init()
	t = 0
end

function update()
	t = t + 1
end

function draw()
	-- background pattern
	for x = 0, W - 1 do
		for y = 0, H - 1 do
			local shade = (x + y + t) % 64
			gfx.pset(x, y, 0, shade, 16)
		end
	end

	-- center of screen
	local cx, cy = 8, 8
	local radius = 4
	local angle = t * 0.05 -- rotation speed

	-- circular motion
	local px = math.floor(cx + math.cos(angle) * radius)
	local py = math.floor(cy + math.sin(angle) * radius)

	gfx.sprite(0, px, py)
end
