const http = require('http');
const fs = require('fs').promises;

const user_data = 'data.json';
let blogPosts = [];

async function loadBlogPosts() {
  try {
    const data = await fs.readFile(user_data, 'utf8');
    blogPosts = JSON.parse(data);
  } catch (err) {
    console.error('Error reading blog posts:', err);
  }
}

loadBlogPosts();

function parseBody(req, res, next) {
  let body = '';
  req.on('data', chunk => {
    body += chunk.toString();
  });
  req.on('end', () => {
    req.body = JSON.parse(body);
    next();
  });
}

const server = http.createServer(async (req, res) => {
  if (req.method === 'POST' && req.url === '/users/add') {
    parseBody(req, res, async () => {
      try {
        const newPost = {
          id: Date.now(),
          title: req.body.title,
          content: req.body.content || '', 
        };
        
        blogPosts.push(newPost);
        
        await fs.writeFile(user_data, JSON.stringify(blogPosts), 'utf8');
        
        res.writeHead(201, {'Content-Type': 'application/json'});
        res.end(JSON.stringify(newPost));
      } catch (error) {
        console.error('An error occurred:', error);
        res.writeHead(500, {'Content-Type': 'text/plain'});
        res.end('Internal Server Error');
      }
    });
  } else if (req.method === 'PUT' && req.url.startsWith('/users/update')) {
    parseBody(req, res, async () => {
      try {
        const urlParts = req.url.split('/');
        const id = urlParts[urlParts.length - 1];
        
        // find the post by userid
        const index = blogPosts.findIndex(post => post.id.toString() === id);
        if (index !== -1) {
          blogPosts[index] = {
            ...blogPosts[index],
            ...req.body,
          };

          await fs.writeFile(user_data, JSON.stringify(blogPosts), 'utf8');
          
          res.writeHead(200, {'Content-Type': 'application/json'});
          res.end(JSON.stringify(blogPosts[index]));
        } else {
          res.writeHead(404, {'Content-Type': 'text/plain'});
          res.end('Resource not found');
        }
      } catch (error) {
        console.error('An error occurred:', error);
        res.writeHead(500, {'Content-Type': 'text/plain'});
        res.end('Internal Server Error');
      }
    });
  } else if (req.method === 'GET') {
    
  } else {
    res.writeHead(405, {'Content-Type': 'text/plain'});
    res.end('Method Not Allowed');
  }
});

const port = 3000;
server.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
