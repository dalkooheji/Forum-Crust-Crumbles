function toggleLike(postID, logged) {
    console.log(logged)
    if (logged === 'false'){
        window.location.href = '/login'
    }
    
    console.log("toggleLike called for PostID:", postID);
    fetch('/toggle-like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ postID: parseInt(postID, 10) }), // Convert postID to an integer
    })
    
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const likeCountElement = document.getElementById('like-count-' + postID);
                likeCountElement.textContent = data.likeCount;
                const dislikeCountElement = document.getElementById('dislike-count-' + postID);
                dislikeCountElement.textContent = data.dislikeCount;

                const likeIconElement = document.getElementById('like-icon-' + postID);
                // Toggle the 'liked' class based on the likeCount
                if (data.dislikeCount > 0) {
                    dislikeIconElement.classList.add('disliked');
                } else {
                 dislikeIconElement.classList.remove('disliked');
                }
                if (data.likeCount > 0) {
                    likeIconElement.classList.add('liked');
                } else {
                    likeIconElement.classList.remove('liked');
                }
            }
        })
        .catch(error => console.error('Error:', error));
        
}

function toggleDislike(postID, logged) {
    console.log(logged)
    if (logged === 'false'){
        window.location.href = '/login'
    }

    console.log("toggleDislike called for PostID:", postID);
    fetch('/toggle-dislike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ postID: parseInt(postID, 10) }), // Convert postID to an integer
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            const dislikeCountElement = document.getElementById('dislike-count-' + postID);
            dislikeCountElement.textContent = data.dislikeCount;
            const likeCountElement = document.getElementById('like-count-' + postID);
            likeCountElement.textContent = data.likeCount;

            const dislikeIconElement = document.getElementById('dislike-icon-' + postID);
            if (data.dislikeCount > 0) {
                dislikeIconElement.classList.add('disliked');
            } else {
                dislikeIconElement.classList.remove('disliked');
            }
            if (data.likeCount > 0) {
                    likeIconElement.classList.add('liked');
            } else {
                    likeIconElement.classList.remove('liked');
            }
        }
    })
    .catch(error => console.error('Error:', error));
    
}

function filterByCategory(category) {
    fetch(`/posts?category=` + encodeURIComponent(category))
        .then(response => response.text())
        .then(html => {
            document.querySelector('.posts-container').innerHTML = html;
        })
        .catch(error => console.error('Error:', error));
}