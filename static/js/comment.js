
function toggleCommentLike(element, logged) {
    console.log(logged)
    if (logged === "false"){
        window.location.href = '/login'
    }
    const commentID = element.getAttribute('data-comment-id');
    toggleCommentReaction(commentID, true, element);
}

function toggleCommentDislike(element, logged) {
    console.log(logged)
    if (logged === "false"){
        window.location.href = '/login'
    }
    const commentID = element.getAttribute('data-comment-id');
    toggleCommentReaction(commentID, false, element);
}


function toggleCommentReaction(commentID, isLike, element) {
    fetch('/toggle-comment-reaction', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            commentID: commentID,
            isLike: isLike
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            const likeCountElement = element.parentElement.querySelector('.comment-like-count');
            const dislikeCountElement = element.parentElement.querySelector('.comment-dislike-count');
            likeCountElement.textContent = data.likeCount;
            dislikeCountElement.textContent = data.dislikeCount;
        } else {
            console.error('Failed to toggle reaction');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
