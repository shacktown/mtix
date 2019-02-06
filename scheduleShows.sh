# schedule movie shows
peer chaincode invoke  -n mycc -C myc -c '{"Args":["ScheduleShow", "Regal1", "1", "Rambo",  "2019-02-05T13:00:00-05:00", "10.00", "100"]}'