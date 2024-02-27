#!/bin/bash
sleep 3

mongosh --host appmongo1:27017 <<EOF
  var cfg = {
    "_id": "myReplicaSet",
    "version": 1,
    "members": [
      {
        "_id": 0,
        "host": "appmongo1:27017",
        "priority": 2
      },
      {
        "_id": 1,
        "host": "appmongo2:27017",
        "priority": 0
      }
    ]
  };
  rs.initiate(cfg);
EOF