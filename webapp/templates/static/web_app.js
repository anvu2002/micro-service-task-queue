var URL_GO = 'http://localhost:9888';
var URL_GO_IMAGE_STATUS = 'http://localhost:9888/image_status?task_id=';
var URL_GO_KEYWORD_STATUS = 'http://localhost:9888/keyword_status?task_id=';
var URL_GO_TTS_STATUS = 'http://localhost:9888/tts_status?task_id=';

var results = [];
var status_list = [];
var res = '';

jQuery(document).ready(function () {
    $('#row_detail').hide();
    $("#row_results").hide();
    $('#progress-bar').hide();

    // MAIN BUTTONs

    $('#btn-get-video').on('click', function () {
        var prompt = $('#file_prompt').val();
        var file_name = encodeURIComponent(prompt);

        $.ajax({
            url: URL_GO + '/process_doc?file_name=' + file_name,
            type: "post",
            beforeSend: function () {
                results = [];
                status_list = [];
                $("#table_result > tbody").html('');
                $('#row_detail').hide();
                $("#row_results").hide();
                console.log("[ML SERVICE] process_doc");
            },
        }).done(function (jsondata, textStatus, jqXHR) {
            console.log(" JSON --- /process_doc = ", jsondata);
            var task_id = jsondata['go_task_id'];
            var status = jsondata['status'];
            console.log("go_task_id = ", task_id);

            status_list.push(task_id);

            var interval = setInterval(refresh, 1000);
            var fake_progress = 0;

            function refresh() {
                var n_success = 0;
                console.log("status_list.length", status_list.length);
                for (var i = 0; i < status_list.length; i++) {
                    $.ajax({
                        url: URL_GO_KEYWORD_STATUS + status_list[i],
                        success: function (data) {
                            console.log("SUCCEED query ", URL_GO_KEYWORD_STATUS + status_list[i]);
                            console.log("/process_doc data = ", data);
                            var id = status_list[i];
                            var status = data['status'];
                            $('#' + id).html(status);
                            if (status == 'SUCCESS' || status == 'FAILED') {
                                $($('#' + id).siblings()[1]).children().show();
                                n_success++;
                            }

                            $('#progress-bar').css('--p', fake_progress);
                            $('#progress-bar').show();
                            if (fake_progress < 100) {
                                fake_progress += Math.floor(Math.random() * 10) + 1;
                            }
                        },
                        async: false
                    });
                }
                if (n_success == status_list.length) {
                    $('#progress-bar').hide();
                    clearInterval(interval);
                    continueTasks(data['keywords']);
                }
            }

        }).fail(function (jsondata, textStatus, jqXHR) {
            console.log(jsondata);
            $("#row_results").hide();
        });
    });

    function continueTasks(keywords) {
        var imageTaskID, sentenceTaskId;

        $.ajax({
            url: URL_GO + '/get_images?query=' + encodeURIComponent(keywords['keywords']),
            type: 'get'
        }).done(function (data) {
            imageTaskID = data['go_task_id'];
        });

        $.ajax({
            url: URL_GO + '/tts?text=' + encodeURIComponent(keywords['sentences']) + '&save_path=./',
            type: 'get'
        }).done(function (data) {
            sentenceTaskId = data['go_task_id'];
        });

        var interval = setInterval(function () {
            checkTaskStatus(imageTaskID, sentenceTaskId, interval);
        }, 1000);
    }

    function checkTaskStatus(imageTaskID, sentenceTaskId, interval) {
        var imageStatus, sentenceStatus;

        $.ajax({
            url: URL_GO_IMAGE_STATUS + imageTaskID,
            type: 'get',
            success: function (data) {
                imageStatus = data['status'];
            },
            async: false
        });

        $.ajax({
            url: URL_GO_TTS_STATUS + sentenceTaskId,
            type: 'get',
            success: function (data) {
                sentenceStatus = data['status'];
            },
            async: false
        });

        if (imageStatus === 'SUCCESS' && sentenceStatus === 'SUCCESS') {
            clearInterval(interval);
            ffmpegTask();
        }
    }

    function ffmpegTask(){
        $.ajax({
            url: URL_GO + '/get_ffmpeg?query=started',
            type: 'get',
            success: function (data) {
                console.log("[Final Task] ffmpeg Completed", data);
            }
        });
    }
});
