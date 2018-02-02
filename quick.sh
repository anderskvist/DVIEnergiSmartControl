#!/bin/bash

COOKIE=${1}

while read LINE; do

  K=$(echo ${LINE}|grep "temp value"|awk -F'[ "]' '{print $4}')
  V=$(echo ${LINE}|grep "temp value"|awk -F'[>&]' '{print $2}')
  P=$(echo ${LINE}|grep "img src" | awk -F '"' '{print $2}' | basename $(cat /dev/stdin) 2> /dev/null)


  if [ -n "${K}" ]; then
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
	  value9)
	      echo dvienergi.smartcontrol.solar_roof value=${V}
	      ;;
	  *)
	      echo "Unhandled temperature"
	      echo ${K}: ${V} 1>&2
	      ;;
      esac
  fi
  if [ -n "${P}" ]; then
      case "${P}" in
	  A1-1.gif)
	      echo dvienergi.smartcontrol.solar_heating_pump value=0i
	      echo dvienergi.smartcontrol.solar_to_ground value=0i
	      ;;
	  A1-1-4.gif)
	      echo dvienergi.smartcontrol.solar_heating_pump value=1i
	      echo dvienergi.smartcontrol.solar_to_ground value=0i
	      ;;
	  A1-2-4.gif)
	      echo dvienergi.smartcontrol.solar_heating_pump value=1i
	      echo dvienergi.smartcontrol.solar_to_ground value=1i
	      ;;
	  A2-1.gif)
	      # not a pump
	      ;;
	  A3-1.gif)
	      # not a pump
	      ;;
	  A4-1.gif)
	      echo dvienergi.smartcontrol.ground_pump value=0i
	      ;;
	  A4-1-4.gif)
	      echo dvienergi.smartcontrol.ground_pump value=1i
	      ;;
	  A4-6-4.gif)
	      echo dvienergi.smartcontrol.ground_pump value=1i
	      ;;
	  A5-1.gif)
	      echo dvienergi.smartcontrol.compressor value=0i
	      ;;
	  A5-1-4.gif)
	      echo dvienergi.smartcontrol.compressor value=1i
	  ;;
	  A6-0.gif)
	      echo dvienergi.smartcontrol.electric_heater value=0i
	      ;;
	  A6-1.gif)
	      echo dvienergi.smartcontrol.electric_heater value=0i
	      ;;
	  A6-1-4.gif)
	      echo dvienergi.smartcontrol.electric_heater value=1i
	      ;;
	  A7-1.gif)
	      echo dvienergi.smartcontrol.house_heating value=0i
	      ;;
	  A7-1-4.gif)
	      echo dvienergi.smartcontrol.house_heating value=1i
	      ;;
	  A8-1.gif)
	      echo dvienergi.smartcontrol.house_heating_pump value=0i
	      ;;
	  A8-1-4.gif)
	      echo dvienergi.smartcontrol.house_heating_pump value=1i
	      ;;

	  *)
	      echo "Unhandled pump"
	      echo ${P} 1>&2
      esac
  fi
done < <(curl 'https://smartcontrol.dvienergi.com/includes/process.php' -H "cookie: PHPSESSID=${COOKIE}" --data 'subupdatepumpgraphics=1' -s|sed 's/></>\n</g'|grep "\(temp value\|img src\)")


while read LINE; do
    case "${LINE}" in
	*user2*)
	    echo ${LINE}|awk -F' ' '{print "dvienergi.smartcontrol.heating_curve value=" $4}'
	    ;;
	*Beregnet\ temperatur*)
	    echo ${LINE}|awk -F'[ <]' '{print "dvienergi.smartcontrol.calculated_setpoint value=" $6}'
	    ;;
    esac
done < <(curl 'https://smartcontrol.dvienergi.com/includes/pumpchoice.php?id=12' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g')

HOTWATERSETPOINT=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpchoice.php?id=22' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep 'user11'|awk -F' ' '{print $4}')



COMPRESSORTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Kompressor" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')
HOTWATERTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Varmt vand" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')
ELECTRICHEATERTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Tilskudsvarme" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')
SOLARTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Solvarme" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')
SOLARTOGROUNDTIME=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpinfo.php?id=31' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep "Sol til Jord" -A 1|tail -n 1|sed -e 's/<[^>]*>//g')

echo dvienergi.smartcontrol.hotwater_setpoint value=${HOTWATERSETPOINT}

echo dvienergi.smartcontrol.compressor_time value=${COMPRESSORTIME}
echo dvienergi.smartcontrol.hotwater_time value=${HOTWATERTIME}
echo dvienergi.smartcontrol.electricheater_time value=${ELECTRICHEATERTIME}
echo dvienergi.smartcontrol.solar_time value=${SOLARTIME}
echo dvienergi.smartcontrol.solartoground_time value=${SOLARTOGROUNDTIME}
