package main

// HTML template with Lightbox
const htmlTemplateStr = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; color: #333; }
        h1 { text-align: center; margin-bottom: 30px; }
        .download-link { display: block; text-align: center; margin-bottom: 40px; }
        .download-link a { background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; transition: background 0.2s; }
        .download-link a:hover { background: #0056b3; }
        
        .gallery { 
            column-width: 450px;
            column-gap: 10px;
            max-width: 1800px; 
            margin: 0 auto; 
        }
        .gallery-item { 
            break-inside: avoid; 
            margin-bottom: 15px; 
        }
        .gallery img { 
            width: 100%; 
            height: auto; 
            border-radius: 2px; 
            display: block; 
            cursor: zoom-in; 
            //box-shadow: 0 2px 5px rgba(0,0,0,0.1); 
            transition: transform 0.2s;
        }
        //.gallery img:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.15); }

        /* Lightbox */
        .lightbox {
            display: none;
            position: fixed;
            z-index: 1000;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.9);
            justify-content: center;
            align-items: center;
        }
        .lightbox.active { display: flex; }
        .lightbox img {
            max-width: 90%;
            max-height: 90vh;
            border-radius: 4px;
            box-shadow: 0 0 20px rgba(0,0,0,0.5);
        }
        .lightbox-close {
            position: absolute;
            top: 20px;
            right: 30px;
            color: white;
            font-size: 40px;
            cursor: pointer;
            user-select: none;
        }
        .lightbox-nav {
            position: absolute;
            top: 50%;
            transform: translateY(-50%);
            color: white;
            font-size: 50px;
            cursor: pointer;
            user-select: none;
            padding: 20px;
            background: rgba(0,0,0,0.2);
            border-radius: 5px;
            transition: background 0.2s;
        }
        .lightbox-nav:hover { background: rgba(0,0,0,0.5); }
        .lightbox-prev { left: 20px; }
        .lightbox-next { right: 20px; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{if .ZipName}}
    <div class="download-link">
        <a href="{{.ZipName}}">Download all photos as ZIP</a>
    </div>
    {{end}}
    
    <section class="gallery">
        {{range .Images}}
        <div class="gallery-item">
            <img class="zoomable" src="thumbnails/{{.}}" data-full="{{$.FolderName}}/{{.}}" alt="{{.}}" loading="lazy">
        </div>
        {{end}}
    </section>

    <div id="lightbox" class="lightbox">
        <span class="lightbox-close">&times;</span>
        <a class="lightbox-nav lightbox-prev">&#10094;</a>
        <a class="lightbox-nav lightbox-next">&#10095;</a>
        <img id="lightbox-img" src="" alt="Lightbox Image">
    </div>

    <script>
        const lightbox = document.getElementById('lightbox');
        const lightboxImg = document.getElementById('lightbox-img');
        const closeBtn = document.querySelector('.lightbox-close');
        const prevBtn = document.querySelector('.lightbox-prev');
        const nextBtn = document.querySelector('.lightbox-next');
        
        let images = [];
        let currentIndex = 0;

        // Collect all zoomable images
        document.querySelectorAll(".zoomable").forEach((img, index) => {
            images.push(img.dataset.full);
            img.addEventListener("click", () => {
                currentIndex = index;
                showImage(currentIndex);
                lightbox.classList.add('active');
            });
        });

        function showImage(index) {
            if (index < 0) index = images.length - 1;
            if (index >= images.length) index = 0;
            currentIndex = index;
            lightboxImg.src = images[currentIndex];
        }

        closeBtn.addEventListener('click', () => {
            lightbox.classList.remove('active');
        });

        lightbox.addEventListener('click', (e) => {
            if (e.target === lightbox) {
                lightbox.classList.remove('active');
            }
        });

        prevBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            showImage(currentIndex - 1);
        });

        nextBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            showImage(currentIndex + 1);
        });

        document.addEventListener('keydown', (e) => {
            if (!lightbox.classList.contains('active')) return;
            
            if (e.key === 'Escape') {
                lightbox.classList.remove('active');
            } else if (e.key === 'ArrowLeft') {
                showImage(currentIndex - 1);
            } else if (e.key === 'ArrowRight') {
                showImage(currentIndex + 1);
            }
        });
    </script>
</body>
</html>`
