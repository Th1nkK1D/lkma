# LK's Neighbor Mimic Matting

## Main algorithm
1. **Read Inputs** - Input image and Trimap Scribbble
2. **Extract Scribble** - Extract FG, BG and Alpha from inputs
3. **Explore Neighbour** - For each unknown pixel, find nearest FG and BG pixel using recursive walk
4. **Mimic Neighbour** - Assume that FG and BG of unknown pixel to have the same value with nearest known FG and GB pixel
5. **Minimize Energy** - Using Gibb's sampling to solve alpha with the mimimum energy
