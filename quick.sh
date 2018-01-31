#!/bin/bash

COOKIE=${1}

while read LINE; do

  K=$(echo ${LINE}|awk -F'[ "]' '{print $4}')
  V=$(echo ${LINE}|awk -F'[>&]' '{print $2}')


  case "${K}" in
    value7)
      echo dvienergi.smartcontrol.outside value=${V}
      ;;
    value5)
      echo dvienergi.smartcontrol.heating value=${V}
      ;;
    value3)
      echo dvienergi.smartcontrol.hotwater value=${V}
      ;;
    value1)
      echo dvienergi.smartcontrol.heating_forward value=${V}
      ;;
    value2)
      echo dvienergi.smartcontrol.heating_return1 value=${V}
      ;;
    value8)
      echo dvienergi.smartcontrol.heating_return2 value=${V}
      ;;
    value13)
      echo dvienergi.smartcontrol.ground_return value=${V}
      ;;
    value14)
      echo dvienergi.smartcontrol.ground_forward value=${V}
      ;;
    value12)
      echo dvienergi.smartcontrol.compressor_front value=${V}
      ;;
    value11)
      echo dvienergi.smartcontrol.compressor_back value=${V}
      ;;
    *)
      echo ${LINE} 1>&2
      echo ${K}: ${V} 1>&2
      ;;

  esac
done < <(curl 'https://smartcontrol.dvienergi.com/includes/process.php' -H "cookie: PHPSESSID=${COOKIE}" --data 'subupdatepumpgraphics=1' -s|sed 's/></>\n</g'|grep "temp value")


HEATINGCURVE=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpchoice.php?id=12' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep 'user2'|awk -F' ' '{print $4}')

CALCULATEDSETPOINT=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpchoice.php?id=12' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep 'Beregnet temperatur'|awk -F'[ <]' '{print $7}')

HOTWATERSETPOINT=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpchoice.php?id=22' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep 'user11'|awk -F' ' '{print $4}')

COMPRESSORTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Kompressor" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')

HOTWATERTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Kompressor" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')

ELECTRICHEATERTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Kompressor" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')

echo dvienergi.smartcontrol.heating_curve value=${HEATINGCURVE}
echo dvienergi.smartcontrol.calculated_setpoint value=${CALCULATEDSETPOINT}
echo dvienergi.smartcontrol.hotwater_setpoint value=${HOTWATERSETPOINT}

echo dvienergi.smartcontrol.compressor_time value=${COMPRESSORTIME}
echo dvienergi.smartcontrol.hotwater_time value=${HOTWATERTIME}
echo dvienergi.smartcontrol.electricheater_time value=${ELECTRICHEATERTIME}
