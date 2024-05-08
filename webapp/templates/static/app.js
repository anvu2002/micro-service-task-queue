// var URL = 'http://localhost:8000'  -- VLNML API endpoints
var URL = 'http://localhost:9888'
var URL_VLNML = 'http://localhost:8000'

var URL_STATUS = 'http://localhost:8000/api/status/'
var results = []
var status_list = []
var res = ''
jQuery(document).ready(function () {
    $('#row_detail').hide()
    $("#row_results").hide();
    $('#progress-bar').hide();

    $('#btn-detect').on('click', function () {
        var form_data = new FormData();
        files = $('#input_file').prop('files')
        for (i = 0; i < files.length; i++)
            form_data.append('files', $('#input_file').prop('files')[i]);

        $.ajax({
            url: URL_VLNML + '/api/process',
            type: "post",
            data: form_data,
            enctype: 'multipart/form-data',
            contentType: false,
            processData: false,
            cache: false,
            beforeSend: function () {
                results = []
                status_list = []
                $("#table_result > tbody").html('');
                $('#row_detail').hide();
                $("#row_results").hide();

            },
        }).done(function (jsondata, textStatus, jqXHR) {
            for (i = 0; i < jsondata.length; i++) {
                task_id = jsondata[i]['task_id']
                status = jsondata[i]['status']
                results.push(URL_VLNML + jsondata[i]['url_result'])
                status_list.push(task_id)
                result_button = `<button class="btn btn-small btn-success" style="display: none" id="btn-view-detect" data=${i}>View</a>`
                $("#table_result > tbody").append(`<tr><td>${task_id}</td><td id=${task_id}>${status}</td><td>${result_button}</td></tr>`);
                $("#row_results").show();
            }

            var interval = setInterval(refresh, 1000);

            var fake_progress = 0
            function refresh() {
                n_success = 0
                for (i = 0; i < status_list.length; i++) {
                    $.ajax({
                        url: URL_STATUS + status_list[i],
                        success: function (data) {
                            id = status_list[i]
                            status = data['status']
                            $('#' + id).html(status)
                            if ((status == 'SUCCESS') || (status == 'FAILED')) {
                                $($('#' + id).siblings()[1]).children().show()
                                n_success++
                            }

                            // Show progress bar and label
                            console.log("fake_progress",fake_progress)
                            $('#progress-bar').css('--p', fake_progress);
                            $('#progress-bar').show();
                            $('#progress-label').text('Processing...');
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
                }
            }
        }).fail(function (jsondata, textStatus, jqXHR) {
            console.log(jsondata)
            $("#row_results").hide();
        });

    })

    $('#btn-caption').on('click', function () {
        var prompt = $('#input_prompt').val();

        console.log(prompt)
        $.ajax({
            url: URL + '/api/get_similarity',
            type: "post",
            data: JSON.stringify({ images: ["./api/images/0.jpg", "./api/images/1.jpg","./api/images/2.jpg","./api/images/3.jpg","./api/images/4.jpg"], prompt: prompt }),
            contentType: 'application/json',
            beforeSend: function () {
                results = []
                status_list = []
                $("#table_result > tbody").html('');
                $('#row_detail').hide();
                $("#row_results").hide();

            },
            success: function (jsondata, textStatus, jqXHR) {
                console.log(jsondata);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.log(jqXHR);
            },
        }).done(function (jsondata, textStatus, jqXHR) {
            for (i = 0; i < jsondata.length; i++) {
                task_id = jsondata[i]['task_id']
                status = jsondata[i]['status']
                results.push(URL + jsondata[i]['url_result'])
                status_list.push(task_id)
                result_button = `<button class="btn btn-small btn-success" style="display: none" id="btn-view-caption" data=${i}>View</a>`
                $("#table_result > tbody").append(`<tr><td>${task_id}</td><td id=${task_id}>${status}</td><td>${result_button}</td></tr>`);
                $("#row_results").show();
            }

            var interval = setInterval(refresh, 1000);
            var fake_progress = 0
            function refresh() {
                n_success = 0
                for (i = 0; i < status_list.length; i++) {
                    $.ajax({
                        url: URL_STATUS + status_list[i],
                        success: function (data) {
                            id = status_list[i]
                            status = data['status']
                            $('#' + id).html(status)
                            if ((status == 'SUCCESS') || (status == 'FAILED')) {
                                $($('#' + id).siblings()[1]).children().show()
                                n_success++
                            }

                            // Show progress bar and label
                            console.log("fake_progress",fake_progress)
                            $('#progress-bar').css('--p', fake_progress);
                            $('#progress-bar').show();
                            $('#progress-label').text('Processing...');
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
                }
            }
        }).fail(function (jsondata, textStatus, jqXHR) {
            console.log(jsondata)
            $("#row_results").hide();
        });
    })

    $('#btn-get-images').on('click', function () {
        var prompt = $('#keyword_prompt').val();
        var query = encodeURIComponent(prompt);

        $.ajax({
            url: URL + '/get_images?query=' + query,
            type: "post",
            beforeSend: function () {
                results = []
                status_list = []
                $("#table_result > tbody").html('');
                $('#row_detail').hide();
                $("#row_results").hide();

            },
            success: function (jsondata, textStatus, jqXHR) {
                console.log(jsondata);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.log(jqXHR);
            },
        }).done(function (jsondata, textStatus, jqXHR) {
            for (i = 0; i < jsondata.length; i++) {
                task_id = jsondata[i]['task_id']
                status = jsondata[i]['status']
                results.push(URL + jsondata[i]['url_result'])
                status_list.push(task_id)
                result_button = `<button class="btn btn-small btn-success" style="display: none" id="btn-view-caption" data=${i}>View</a>`
                $("#table_result > tbody").append(`<tr><td>${task_id}</td><td id=${task_id}>${status}</td><td>${result_button}</td></tr>`);
                $("#row_results").show();
            }

            var interval = setInterval(refresh, 1000);
            var fake_progress = 0
            function refresh() {
                n_success = 0
                for (i = 0; i < status_list.length; i++) {
                    $.ajax({
                        url: URL_STATUS + status_list[i],
                        success: function (data) {
                            id = status_list[i]
                            status = data['status']
                            $('#' + id).html(status)
                            if ((status == 'SUCCESS') || (status == 'FAILED')) {
                                $($('#' + id).siblings()[1]).children().show()
                                n_success++
                            }

                            // Show progress bar and label
                            console.log("fake_progress",fake_progress)
                            $('#progress-bar').css('--p', fake_progress);
                            $('#progress-bar').show();
                            $('#progress-label').text('Processing...');
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
                }
            }
        }).fail(function (jsondata, textStatus, jqXHR) {
            console.log(jsondata)
            $("#row_results").hide();
        });
    })

    $(document).on('click', '#btn-view-detect', function (e) {
        id = $(e.target).attr('data')
        $.get(results[id], function (data) {
            if (data) {
                $('#row_detail').show()
                $('#result_txt').val(JSON.stringify(data['bbox'], undefined, 4))
                $('#result_img').attr('src', URL_VLNML + '/' + data.file_name)
                $('#result_link').attr('href', URL_VLNML + '/' + data.file_name)
            } else {
                alert('Result not ready or already consumed!')
                $('#row_detail').hide()
            }
        });
    })

    $(document).on('click', '#btn-view-caption', function (e) {
        id = $(e.target).attr('data')
        $.get(results[id], function (data) {
            res = data
            if (data['status'] == 'SUCCESS') {
                for (var i = 0; i < res.result.length; i++) {
                    $('#row_detail').show()
                    $('#result_txt').val(JSON.stringify(res.result[i], undefined, 4))
                    $('#result_img').attr('src', `${URL}/img/${i}.jpg`)
                    $('#result_link').attr('href', `${URL}/img/${i}.jpg`)
                }
            } else {
                alert('Result not ready or already consumed!')
                $('#row_detail').hide()
            }
        });
    })


    $(document).on('click', '#btn-refresh', function (e) {
        for (i = 0; i < status_list.length; i++) {
            $.ajax({
                url: URL_STATUS + status_list[i],
                success: function (data) {
                    id = status_list[i]
                    status = data['status']
                    $('#' + id).html(status)
                    if (status == 'SUCCESS')
                        $($('#' + id).siblings()[1]).children().show()
                },
                async: false
            });
        }
    })

})
