#pragma semicolon 1

#include <sourcemod>
#include <Source-Chat-Relay>
#include <ccc>

char g_sChannelname[64];
ConVar g_cChannelname;
EngineVersion g_GameEngine = Engine_Unknown;

public Plugin myinfo = 
{
	name = "Admin Relay -SCR",
	author = "granny",
	description = "Admin chat extension for Source Chat Relay",
	version = "1.0",
	url = "https://github.com/pkrok01/source-chat-relay/tree/personal"
};

public void OnPluginStart()
{
	g_GameEngine = GetEngineVersion();
	
	g_cChannelname = CreateConVar("g_ac_chatname", "", "The name of the channel you would like to relay the admin chat to (Make sure to capitalize each word)", FCVAR_NONE);

	AutoExecConfig(true);
}

public void OnConfigsExecuted()
{
	g_cChannelname.GetString(g_sChannelname, sizeof g_sChannelname);
}

public Action SCR_OnMessageReceive(const char[] sEntityName, IdentificationType iIDType, const char[] sID, char[] sClientName, char[] sMessage)
{
	g_cChannelname.GetString(g_sChannelname, sizeof g_sChannelname);

	if (StrEqual(g_sChannelname, sEntityName))
	{
		for (int i = 1; i <= MaxClients; i++)
		{
			if (IsClientInGame(i) && (CheckCommandAccess(i, "sm_chat", ADMFLAG_CHAT)))
			{
				if (g_GameEngine == Engine_CSGO)
					PrintToChat(i, " \x01\x0B\x04(ADMINS) %s: \x01%s", sClientName, sMessage);
				else
					PrintToChat(i, "\x04(ADMINS) %s: \x01%s", sClientName, sMessage);
			}	
		}
		return Plugin_Stop;
	}
	return Plugin_Continue;
}