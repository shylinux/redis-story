Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) { can.require(["/plugin/story/studiolayout.js"], function() {
		can.onimport.project(can, msg, aaa.SESS, function(event, sess, value) { var ui = {}; return {
			display: {index: "web.code.redis.shells", args: sess, _init: function(sub) { ui.display = sub }},
			content: {index: "web.code.redis.keys", args: sess},
			profile: can.isCmdMode() && can.onimport._commands(can, sess, ui),
		} })
	}) },
	_commands: function(can, sess, ui) {
		return {index: "web.code.redis.commands", args: sess, _init: function(sub) {
			sub.onexport.output = function(sub) {
				sub.onaction.helpCmd = function(event, can) { var xterm = ui.display._plugins[0].sub, target = event.target
					var msg = can.request(event)
					function input(cmd) { msg.detail = ["input", "only", cmd+"\r\n"]; xterm.onimport.input(xterm, msg) }
					input("clear")
					can.onmotion.delay(can, function() { var begin = xterm.current._buffer.normal.cursorY
						input("help"+" "+msg.Option("command"))
						can.onmotion.delay(can, function() { var end = xterm.current._buffer.normal.cursorY;
							xterm.current.selectLines(begin+1, end-1)
							var ls = xterm.current.getSelection().trim().split("\n")
							ls && ls.length > 3 && sub.Update(sub.request({target: target}, {
								command: msg.Option("command"),
								type: can.base.trimPrefix(ls[3].trim(), "group: "),
								name: can.base.trimPrefix(ls[0].trim(), msg.Option("command").toUpperCase()+" "),
								text: can.base.trimPrefix(ls[1].trim(), "summary: "),
							}), [ctx.ACTION, mdb.CREATE])
						}, 300)
					}, 300)
				}
			}
		}}
	},
	_nick: function(can, value) { return value.sess.slice(0, 6)+`(${value.host}:${value.port}) ${value.role}` },
}, [""])
