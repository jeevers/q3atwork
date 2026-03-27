#!/bin/sh
echo "Edit this script to change the path to ioquake3's dedicated server executable"
echo "Set the sv_dlURL setting to a url like http://yoursite.com/ioquake3_path for ioquake3 clients to download extra data."

# sv_dlURL needs to have quotes escaped like \"http://yoursite.com/ioquake3_path\" or it will be set to "http:" in-game.
IOQUAKE_ROOT=/usr/local/ioquake3

exec_config="+exec extra_configs.cfg"
starting_map="test_bigbox"
if [ -z "$IOQ3_RCON_PASSWORD" ]
then
  echo "rcon password not set! exiting..."
  exit 78 
else
  echo "seta rconPassword \"${IOQ3_RCON_PASSWORD}\"
  map ${starting_map}
  " > ${IOQUAKE_ROOT}/baseq3/extra_configs.cfg
fi

${IOQUAKE_ROOT}/ioq3ded.x86_64 +set dedicated 1 +set sv_allowDownload 1 +set sv_dlURL \"\" +set com_hunkmegs 64 $exec_config "$@"
