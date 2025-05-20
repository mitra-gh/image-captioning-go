const express = require("express");
const path = require("path");

const app = express();
const PORT = 5500; // Or any other port you like

// Serve static files from the "public" folder
app.use(express.static(path.join(__dirname, "public")));

// Optional: serve index.html directly if needed
app.get("/", (req, res) => {
  res.sendFile(path.join(__dirname, "public", "index.html"));
});

app.listen(PORT, () => {
  console.log(`🚀 Server is running at http://localhost:${PORT}`);
});
