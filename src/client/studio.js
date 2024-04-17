Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) { can.require(["/plugin/story/studiolayout.js"], function() {
		can.onimport.project(can, msg, aaa.SESS, function(event, sess, value) { return {
			content: {index: "web.code.redis.keys", args: sess},
			profile: {index: "web.code.redis.commands", args: sess},
			display: {index: "web.code.redis.shells", args: sess},
		} })
	}) },
	_nick: function(can, value) {
		return value.sess.slice(0, 6)+`(${value.host}:${value.port}) ${value.role}`
	},
}, [""])
