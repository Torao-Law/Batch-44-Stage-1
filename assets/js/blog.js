// global array
let blogs = []

// function get data from input text html
function getData(e) {
    e.preventDefault()

    // delaclaration variable dom selection
    let title = document.getElementById('input-blog-title').value
    let content = document.getElementById('input-blog-content').value
    let image = document.getElementById('input-blog-image').files

    // Convert spesific image to blob object
    image = URL.createObjectURL(image[0])

    let dataBlog = {
        title,
        content,
        image,
        author: "Dandi Saputra",
        postedAt: new Date()
    }
    blogs.push(dataBlog)
    console.log(blogs)
    showData()
}

// declaration function show list preview data blog
function showData(){
    document.getElementById("contents").innerHTML = ""
    for(let i=0; i<= blogs.length; i++){
        document.getElementById("contents").innerHTML += `
        <div class="blog-list-item">
            <div class="blog-image">
                <img src="${blogs[i].image}" alt="" />
            </div>
            <div class="blog-content">
            <div class="btn-group">
                <button class="btn-edit">Edit Post</button>
                <button class="btn-post">Post Blog</button>
            </div>
            <h1>
                <a href="blog-detail.html" target="_blank">${blogs[i].title}</a>
            </h1>
            <div class="detail-blog-content">
                ${blogs[i].postedAt} | ${blogs[i].author}
            </div>
            <p>${blogs[i].content}</p>

            <div style="float:right; margin: 10px">
                <p style="font-size: 15px; color:grey">1 minutes ago</p>
            </div>

            </div>
        </div>
        `
    }
}

