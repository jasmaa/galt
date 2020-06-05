# Project

## API
  - /user
    - /:userID : GET
    - /new : POST (auth)
    - /edit : PUT (auth)
    - /delete : DELETE (auth)
    - /:userID/statuses : GET
    - /:userID/circles : GET
    - /:userID/circles/:circleID/new : POST (auth)
    - /:userID/circles/:circleID/edit : PUT (auth)
    - /:userID/circles/:circleID/addUser : POST (auth)
    - /:userID/circles/:circleID/removeUser : POST (auth)
    - /:userID/circles/:circleID/delete : DELETE (auth)

  - /status
    - /new : POST (auth)
    - /:statusID : GET
    - /:statusID/edit : PUT (auth)
    - /:statusID/delete : DELETE (auth)  

  - /group
    - /:groupID : GET
    - /new : POST (auth)
    - /:groupID/edit : PUT (auth)
    - /:groupID/delete : DELETE (auth)
    - /:groupID/statuses : GET (auth)