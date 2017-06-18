package main

import "fmt"
import "bufio"
import "os"
import "path/filepath"

import "github.com/codecat/go-libs/log"

func (self *MPServer) WriteConfig(fnm string) bool {
	targetDir := filepath.Dir(fnm)
	if !pathExists(targetDir) {
		os.MkdirAll(targetDir, os.ModePerm)
	}

	out, err := os.Create(fnm)
	if err != nil {
		log.Fatal("Couldn't create config file \"%s\": %s", fnm, err.Error())
		return false
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	fmt.Fprintf(w, "<!-- Auto-generated with mpman -->\n")
	fmt.Fprintf(w, "<dedicated>\n")

	// XMLRPC passwords
	fmt.Fprintf(w, "  <authorization_levels>\n")
	fmt.Fprintf(w, "    <level>\n")
	fmt.Fprintf(w, "      <name>SuperAdmin</name>\n")
	fmt.Fprintf(w, "      <password>%s</password>\n", self.Passwords.SuperAdmin)
	fmt.Fprintf(w, "    </level>\n")
	fmt.Fprintf(w, "    <level>\n")
	fmt.Fprintf(w, "      <name>Admin</name>\n")
	fmt.Fprintf(w, "      <password>%s</password>\n", self.Passwords.Admin)
	fmt.Fprintf(w, "    </level>\n")
	fmt.Fprintf(w, "    <level>\n")
	fmt.Fprintf(w, "      <name>User</name>\n")
	fmt.Fprintf(w, "      <password>User</password>\n")
	fmt.Fprintf(w, "    </level>\n")
	fmt.Fprintf(w, "  </authorization_levels>\n")

	// Masterserver
	fmt.Fprintf(w, "  <masterserver_account>\n")
	fmt.Fprintf(w, "    <login>%s</login>\n", self.Info["Login"].(string))
	fmt.Fprintf(w, "    <password>%s</password>\n", self.Info["LoginPassword"].(string))
	fmt.Fprintf(w, "    <validation_key>%s</validation_key>\n", self.Info["LoginKey"].(string))
	fmt.Fprintf(w, "  </masterserver_account>\n")

	// Server options
	fmt.Fprintf(w, "  <server_options>\n")
	fmt.Fprintf(w, "    <name>%s</name>\n", self.Info["Name"].(string))
	fmt.Fprintf(w, "    <comment>Managed by $l[nimz.se]Nimz.se</comment>\n")
	fmt.Fprintf(w, "    <hide_server>0</hide_server>\n")

	fmt.Fprintf(w, "    <max_players>%d</max_players>\n", self.Info["MaxPlayers"].(int))
	fmt.Fprintf(w, "    <password_spectator>%s</password_spectator>\n", self.Info["SpectatePassword"].(string))

	fmt.Fprintf(w, "    <max_spectators>%d</max_spectators>\n", self.Info["MaxPlayers"].(int))
	fmt.Fprintf(w, "    <password>%s</password>\n", self.Info["ConnectPassword"].(string))

	fmt.Fprintf(w, "    <keep_player_slots>False</keep_player_slots>\n")
	fmt.Fprintf(w, "    <ladder_mode>forced</ladder_mode>\n")
	fmt.Fprintf(w, "    <enable_p2p_upload>True</enable_p2p_upload>\n")
	fmt.Fprintf(w, "    <enable_p2p_download>False</enable_p2p_download>\n")

	fmt.Fprintf(w, "    <callvote_timeout>60000</callvote_timeout>\n")
	fmt.Fprintf(w, "    <callvote_ratio>0.5</callvote_ratio>\n")
	fmt.Fprintf(w, "    <callvote_ratios>\n")
	fmt.Fprintf(w, "      <voteratio command=\"Ban\" ratio=\"-1\" />\n")
	fmt.Fprintf(w, "    </callvote_ratios>\n")

	fmt.Fprintf(w, "    <allow_map_download>True</allow_map_download>\n")
	fmt.Fprintf(w, "    <autosave_replays>False</autosave_replays>\n")
	fmt.Fprintf(w, "    <autosave_validation_replays>False</autosave_validation_replays>\n")

	fmt.Fprintf(w, "    <referee_password></referee_password>\n")
	fmt.Fprintf(w, "    <referee_validation_mode>0</referee_validation_mode>\n")
	fmt.Fprintf(w, "    <use_changing_validation_seed>False</use_changing_validation_seed>\n")

	fmt.Fprintf(w, "    <disable_horns>False</disable_horns>\n")
	fmt.Fprintf(w, "    <clientinputs_maxlatency>0</clientinputs_maxlatency>\n")
	fmt.Fprintf(w, "  </server_options>\n")

	// System config
	fmt.Fprintf(w, "  <system_config>\n")
	fmt.Fprintf(w, "    <title>%s</title>\n", self.Info["Title"].(string))

	fmt.Fprintf(w, "    <bind_ip_address></bind_ip_address>\n")
	fmt.Fprintf(w, "    <server_port>%d</server_port>\n", self.Info["Port"].(int))
	fmt.Fprintf(w, "    <server_p2p_port>%d</server_p2p_port>\n", self.Info["PortP2P"].(int))
	fmt.Fprintf(w, "    <xmlrpc_port>%d</xmlrpc_port>\n", self.Info["PortRPC"].(int))
	fmt.Fprintf(w, "    <xmlrpc_allowremote>False</xmlrpc_allowremote>\n")

	fmt.Fprintf(w, "    <connection_uploadrate>8000</connection_uploadrate>\n")
	fmt.Fprintf(w, "    <connection_downloadrate>8000</connection_downloadrate>\n")
	fmt.Fprintf(w, "    <allow_spectator_relays>False</allow_spectator_relays>\n")
	fmt.Fprintf(w, "    <p2p_cache_size>600</p2p_cache_size>\n")

	fmt.Fprintf(w, "    <force_ip_address></force_ip_address>\n")
	fmt.Fprintf(w, "    <client_port>0</client_port>\n")
	fmt.Fprintf(w, "    <use_nat_upnp></use_nat_upnp>\n")
	fmt.Fprintf(w, "    <gsp_name>Nimz.se</gsp_name>\n")
	fmt.Fprintf(w, "    <gsp_url>https://nimz.se/</gsp_url>\n")
	fmt.Fprintf(w, "    <scriptcloud_source>nadeocloud</scriptcloud_source>\n")

	fmt.Fprintf(w, "    <blacklist_url></blacklist_url>\n")
	fmt.Fprintf(w, "    <guestlist_filename></guestlist_filename>\n")
	fmt.Fprintf(w, "    <blacklist_filename></blacklist_filename>\n")

	fmt.Fprintf(w, "    <minimum_client_build></minimum_client_build>\n")
	fmt.Fprintf(w, "    <disable_coherence_checks>False</disable_coherence_checks>\n")
	fmt.Fprintf(w, "    <disable_replay_recording>False</disable_replay_recording>\n")
	fmt.Fprintf(w, "    <use_proxy>False</use_proxy>\n")
	fmt.Fprintf(w, "    <proxy_login></proxy_login>\n")
	fmt.Fprintf(w, "    <proxy_password></proxy_password>\n")
	fmt.Fprintf(w, "  </system_config>\n")

	fmt.Fprintf(w, "</dedicated>\n")

	return true
}

func (self *MPServer) WriteMatchSettings(fnm string) bool {
	targetDir := filepath.Dir(fnm)
	if !pathExists(targetDir) {
		os.MkdirAll(targetDir, os.ModePerm)
	}

	out, err := os.Create(fnm)
	if err != nil {
		log.Fatal("Couldn't create match settings file \"%s\": %s", fnm, err.Error())
		return false
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	fmt.Fprintf(w, "<playlist>\n")

	// Game info
	fmt.Fprintf(w, "  <gameinfos>\n")
	fmt.Fprintf(w, "    <game_mode>0</game_mode>\n")
	fmt.Fprintf(w, "    <chat_time>10000</chat_time>\n")
	fmt.Fprintf(w, "    <finishtimeout>1</finishtimeout>\n")
	fmt.Fprintf(w, "    <allwarmupduration>0</allwarmupduration>\n")
	fmt.Fprintf(w, "    <disablerespawn>0</disablerespawn>\n")
	fmt.Fprintf(w, "    <forceshowallopponents>0</forceshowallopponents>\n")
	fmt.Fprintf(w, "    <script_name><![CDATA[TimeAttack.Script.txt]]></script_name>\n")
	fmt.Fprintf(w, "  </gameinfos>\n")

	// Filter
	fmt.Fprintf(w, "  <filter>\n")
	fmt.Fprintf(w, "    <is_lan>1</is_lan>\n")
	fmt.Fprintf(w, "    <is_internet>1</is_internet>\n")
	fmt.Fprintf(w, "    <is_solo>0</is_solo>\n")
	fmt.Fprintf(w, "    <is_hotseat>0</is_hotseat>\n")
	fmt.Fprintf(w, "    <sort_index>1000</sort_index>\n")
	fmt.Fprintf(w, "    <random_map_order>0</random_map_order>\n")
	fmt.Fprintf(w, "  </filter>\n")

	// Script settings
	fmt.Fprintf(w, "  <script_settings>\n")
	//
	fmt.Fprintf(w, "  </script_settings>\n")

	// Maps
	fmt.Fprintf(w, "  <startindex>0</startindex>\n")

	maps := dbQuery("SELECT * FROM maps WHERE Server=?", self.ID())
	for _, m := range maps {
		fmt.Fprintf(w, "  <map>\n")
		fmt.Fprintf(w, "    <file>%s</file>\n", m["Filename"].(string))
		fmt.Fprintf(w, "  </map>\n")
	}

	fmt.Fprintf(w, "</playlist>\n")

	return true
}
