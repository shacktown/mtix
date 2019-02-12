# schedule movie shows
# For each movie time (3:00, 5:00, 7:00, 9:00)
#   schedule a movie in each hall (1-5)
#   for this example the theater is Regal1
#   an incrementing 'showID' is created as the unique key for each movie showing
set -x
CC=mtix
CHAN=myc
theater="Regal1"
ID=0
for time in 15, 17, 19, 21
do
    for i in {1..5}    
    do
    #peer chaincode invoke  -n mtix -C myc -c '{"Args":["ScheduleShow", "show'$ID'", "'$theater'", "hall'$i'", "Rambo '$i''",  "2019-02-14T13:00:00-05:00", "10.00", "100", "0"]}'
    #echo \"ScheduleShow\", \"show$ID\", \"$theater\", \"hall$i\", \"Rambo $i\",  \"2019-02-14T$time:00:00-05:00\", \"10.00\", \"100\", \"0\"
    peer chaincode invoke  -n $CC -C $CHAN -c '{"Args":["ScheduleShow", "show'$ID'", "'$theater'", "hall'$i'", "Rambo '$i'",  "2019-02-14T'$time':00:00-05:00", "10.00", "100", "0"]}'
 
    ID=$((ID + 1))
    done
done

#peer chaincode invoke  -n mtix -C myc -c '{"Args":["ScheduleShow", "show'$ID'", "'$theater'", "hall'$i'", "Rambo '$i''",  "2019-02-14T13:00:00-05:00", "10.00", "100", "0"]}'
#peer chaincode invoke  -n mycc -C mychannel -c '{"Args":["ScheduleShow", "Regal1", "hall1", "Rambo",  "2019-02-05T13:00:00-05:00", "10.00", "100", "0"]}'
#peer chaincode invoke  -n mycc -C mychannel -c '{"Args":["get", "show1"]}'