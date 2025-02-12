// db.auth('butler', 'GGGGgggHJNYTYRTCKGV__lknbjbG89768645edcy564e')

db = db.getSiblingDB('butler_kitten')

db.createUser({
    user: 'butler_app',
    pwd: 'app_secret_GKGILUG79697698__jhvhl_app_secret',
    roles: [
        {
            role: 'readWrite',
            db: 'butler_kitten'
        }
    ]
})

