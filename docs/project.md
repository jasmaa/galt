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
    - /:userID/circles/:circleID : POST (auth)
    - /:userID/circles/:circleID : PUT (auth)
    - /:userID/circles/:circleID/addUser : POST (auth)
    - /:userID/circles/:circleID/removeUser : POST (auth)
    - /:userID/circles/:circleID : DELETE (auth)

  - /status
    - / : POST (auth)
    - /:statusID : GET
    - /:statusID : PUT (auth)
    - /:statusID : DELETE (auth)  

  - /group
    - /:groupID : GET
    - / : POST (auth)
    - /:groupID : PUT (auth)
    - /:groupID : DELETE (auth)
    - /:groupID/statuses : GET (auth)