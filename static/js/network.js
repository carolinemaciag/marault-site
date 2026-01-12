// Grab the canvas
const canvas = document.getElementById("network-canvas");
const ctx = canvas.getContext("2d");

// Canvas size
let width, height;

// Network settings
const NODE_COUNT = 70;
const MAX_DISTANCE = 130;
let nodes = [];

// Resize canvas to window
function resizeCanvas() {
  width = canvas.width = window.innerWidth;
  height = canvas.height = window.innerHeight;
}

window.addEventListener("resize", resizeCanvas);
resizeCanvas();

// Create nodes
function createNodes() {
  nodes = [];
  for (let i = 0; i < NODE_COUNT; i++) {
    nodes.push({
      x: Math.random() * width,
      y: Math.random() * height,
      vx: (Math.random() - 0.5) * 0.2,
      vy: (Math.random() - 0.5) * 0.2
    });
  }
}

createNodes();

// Animation loop
function animate() {
  ctx.clearRect(0, 0, width, height);

  for (let i = 0; i < nodes.length; i++) {
    const a = nodes[i];

    // Move node
    a.x += a.vx;
    a.y += a.vy;

    // Bounce off edges
    if (a.x <= 0 || a.x >= width) a.vx *= -1;
    if (a.y <= 0 || a.y >= height) a.vy *= -1;

    // Connect nearby nodes
    for (let j = i + 1; j < nodes.length; j++) {
      const b = nodes[j];
      const dx = a.x - b.x;
      const dy = a.y - b.y;
      const dist = Math.sqrt(dx * dx + dy * dy);

      if (dist < MAX_DISTANCE) {
        ctx.strokeStyle = `rgba(201, 162, 77, ${1 - dist / MAX_DISTANCE})`;
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(a.x, a.y);
        ctx.lineTo(b.x, b.y);
        ctx.stroke();
      }
    }

    // Draw node
    ctx.fillStyle = "#c9a24d";
    ctx.beginPath();
    ctx.arc(a.x, a.y, 2, 0, Math.PI * 2);
    ctx.fill();
  }

  requestAnimationFrame(animate);
}

animate();
