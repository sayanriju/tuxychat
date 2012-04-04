var active = true
var toggling = false
var notify = false
var permAsked = false

$(function() {
    resizeChat()
    $('#chatBox').focus()

    notify = !!window.webkitNotifications

    $(window).resize(resizeChat)
    $(window).blur(function() {
        active = false
    })
    $(window).focus(function() {
        active = true
    })

	var ch = new goog.appengine.Channel(token)
	var sock = ch.open()
	sock.onmessage = function(msg) {
		var cmd = JSON.parse(msg.data)
		if(cmd.Type == 'join') {
			$('#chatLog').append('<span class="join">' + cmd.Email + ' has joined</span><br>')
		} else if(cmd.Type == 'msg') {
			$('#chatLog').append('<span class="msg"><span class="email">' + cmd.Email + '</span>: ' + cmd.Message + '</span><br>')
            if(!active) {
                if(notify) {
                    popup = window.webkitNotifications.createNotification("", 'Tuxy Chat', 'New message from ' + cmd.Email)
                    popup.onclick = function() {
                        window.focus()
                        $('#chatBox').focus()
                        this.cancel()
                    }

                    popup.show()
                    setTimeout(function() {
                        popup.cancel()
                    }, 5000)
                } else if(!toggling) {
                    toggling = true
                    toggleTitle()
                }
            }
		}
		scrollToEnd()
	}
	
	$('#chatBox').keypress(function(evt) {
		if(evt.which == 13) {
            if(notify) {
                window.webkitNotifications.requestPermission(function() {
                    notify = window.webkitNotifications.checkPermission() == 0
                })
            }

            if($(this).val() == '') {
                return
            }

			$('#chatLog').append('<span class="myMsg"><span class="email">' + email + '</span>: ' + $(this).val() + '</span><br>')
            scrollToEnd()

			$.ajax('/msg/' + roomId, {
				type: 'POST',
				data: {
					msg: $(this).val()
				},
				error: function() {
					$('.error').text('Publishing chat message failed!')
				}
			})
			$(this).val('')
		}
	})
})

function resizeChat() {
    $('#chatLog').height($(window).innerHeight() - 175)
    scrollToEnd()
}

function scrollToEnd() {
    $('#chatLog').scrollTop(50000)
}

function toggleTitle() {
    var title = document.title
    change()

    function change() {
        document.title = "New Message"
        setTimeout(revert, 500)
    }

    function revert() {
        document.title = title
        if(!active) setTimeout(change, 500)
    }
}
