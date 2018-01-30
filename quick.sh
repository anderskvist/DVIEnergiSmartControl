#!/bin/bash

COOKIE=${1}

while read LINE; do

  K=$(echo ${LINE}|awk -F'[ "]' '{print $4}')
  V=$(echo ${LINE}|awk -F'[>&]' '{print $2}')


  case "${K}" in
    value7)
      echo dvienergi.smartcontrol.outside value=${V}0
      ;;
    value5)
      echo dvienergi.smartcontrol.heating value=${V}0
      ;;
    value3)
      echo dvienergi.smartcontrol.hotwater value=${V}0
      ;;
    value1)
      echo dvienergi.smartcontrol.heating_forward value=${V}0
      ;;
    value2)
      echo dvienergi.smartcontrol.heating_return1 value=${V}0
      ;;
    value8)
      echo dvienergi.smartcontrol.heating_return2 value=${V}0
      ;;
    value13)
      echo dvienergi.smartcontrol.ground_return value=${V}0
      ;;
    value14)
      echo dvienergi.smartcontrol.ground_forward value=${V}0
      ;;
    value12)
      echo dvienergi.smartcontrol.compressor_front value=${V}0
      ;;
    value11)
      echo dvienergi.smartcontrol.compressor_back value=${V}0
      ;;
    *)
      echo ${LINE} 1>&2
      echo ${K}: ${V} 1>&2
      ;;

  esac
done < <(curl 'https://smartcontrol.dvienergi.com/includes/process.php' -H "cookie: PHPSESSID=${COOKIE}" --data 'subupdatepumpgraphics=1' -s|sed 's/></>\n</g'|grep "temp value")


HEATINGCURVE=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpchoice.php?id=12' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep 'user2'|awk -F' ' '{print $4}')

CALCULATEDSETPOINT=$(curl 'https://smartcontrol.dvienergi.com/includes/pumpchoice.php?id=12' -H "cookie: PHPSESSID=${COOKIE}" -s|sed 's/></>\n</g'|grep 'Beregnet temperatur'|awk -F'[ <]' '{print $7}')


echo dvienergi.smartcontrol.heating_curve value=${HEATINGCURVE}
echo dvienergi.smartcontrol.calculated_setpoint value=${CALCULATEDSETPOINT}
