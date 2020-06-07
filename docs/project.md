# Project

## API
  - /user
    - /:userID : GET
    - / : GET (auth)
    - / : POST (auth)
    - / : PUT (auth)
    - / : DELETE (auth)
    - /:userID/statuses : GET
    - /:userID/circles : GET

  - /status
    - /:statusID : GET
    - / : POST (auth)
    - /:statusID : PUT (auth)
    - /:statusID : DELETE (auth)
    - /:statusID/comments : GET
    - /:statusID/comment : POST (auth)
    
  - /circle
    - /:circleID : POST (auth)
    - /:circleID : PUT (auth)
    - /:circleID/addUser : POST (auth)
    - /:circleID/removeUser : POST (auth)
    - /:circleID : DELETE (auth)

  - /comment
    - /:commentID : GET
    - /:commentID : PUT (auth)
    - /:commentID : DELETE (auth)
    - /:commentID/comment : POST (auth)

  - /group
    - /:groupID : GET
    - / : POST (auth)
    - /:groupID : PUT (auth)
    - /:groupID : DELETE (auth)
    - /:groupID/statuses : GET (auth)
    - /:groupID/status : POST (auth)
