// var URL = 'http://localhost:8000'  -- VLNML API endpoints
var URL_GO = 'http://localhost:9888'
var URL_GO_IMAGE_STATUS = 'http://localhost:9888/image_status?task_id='
var URL_GO_KEYWORD_STATUS = 'http://localhost:9888/keyword_status?task_id='

var URL_VLNML = 'http://localhost:8000'

var URL_VLNML_STATUS = 'http://localhost:8000/api/status/'
var results = []
var status_list = []
var res = ''
jQuery(document).ready(function () {
    $('#row_detail').hide()
    $("#row_results").hide();
    $('#progress-bar').hide();

    // MAIN Buttons
    $('#btn-detect').on('click', function () {
        // Phase 1: Request ML Service
        var form_data = new FormData();
        files = $('#input_file').prop('files')
        for (i = 0; i < files.length; i++)
            form_data.append('files', $('#input_file').prop('files')[i]);

        $.ajax({
            url: URL_VLNML + '/api/process_detection',
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
            success: function(jsondata, textStatus, jqXHR) {
                console.log("[SUCCEED] Requested Detect Objects!");
                console.log("JSON DATA = ",jsondata);
                console.log("JSON DATA LEN =", jsondata.length)
                console.log("jqXHR status",jqXHR.status)
            },
            error: function(jqXHR, textStatus, errorThrown){
                console.log("[FAILED] Request Detect Objects.");
                console.log("jqXHR = ",jqXHR)
            }
        }).done(function (jsondata, textStatus, jqXHR) {
            console.log("[SUCCEED] Requested Detect Objects!");
            console.log("jqXHR status",jqXHR.status);
            // jsondata.length == number of submmited pictures (tasks)
            for (i = 0; i < jsondata.length; i++) {
                task_id = jsondata[i]['task_id']
                status = jsondata[i]['status']
                results.push(URL_VLNML + jsondata[i]['url_result'])
                status_list.push(task_id)
                result_button = `<button class="btn btn-small btn-success" style="display: none" id="btn-view-detect" data=${i}>View</a>`
                $("#table_result > tbody").append(`<tr><td>${task_id}</td><td id=${task_id}>${status}</td><td>${result_button}</td></tr>`);
                $("#row_results").show();
            }

            // Phase 2: Probe requested ML Service status
            var interval = setInterval(refresh, 1000);
            var fake_progress = 0;

            function refresh() {
                n_success = 0
                for (i = 0; i < status_list.length; i++) {
                    // Loop through each task
                    $.ajax({
                        url: URL_VLNML_STATUS + status_list[i],
                        success: function (data,textStatus,jqXHR) {
                            id = status_list[i]
                            status = data['status']
                            console.log("[+] Task ",i," -- ",id)
                            console.log("[+] status = ", status)
                            console.log("[+] data = ", data)
                            console.log("[+] jqXHR = ", jqXHR)
                            $('#' + id).html(status)
                            if ((status == 'SUCCESS') || (status == 'FAILED')) {
                                // show View Button for each Completed Task
                                console.log("[+] show View Button")
                                console.log("[+] jqXHR = ", jqXHR)
                                console.log("[+]  $($('#' + id) = ", $($('#' + id)))
                                console.log("[+]  $($('#' + id).siblings()[1]) = ", $($('#' + id)).siblings()[1])
                                console.log("[+]  $($('#' + id).siblings()[1]).children() = ", $($('#' + id).siblings()[1]).children())
                                $($('#' + id).siblings()[1]).children().show()
                                n_success++
                            }

                            // Show progress bar
                            console.log("fake_progress",fake_progress)
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
                    // Complete prob status for ALL tasks
                    $('#progress-bar').hide();
                    clearInterval(interval);
                }
            }
        }).fail(function (jsondata, textStatus, jqXHR) {
            console.log(jsondata)
            $("#row_results").hide();
        });

    })

    $('#btn-get-video').on('click', function () {
        // Phase 1: Request Go Service
        var prompt = $('#file_prompt').val();
        var file_name = encodeURIComponent(prompt);

        $.ajax({
            url: URL_GO + '/process_doc?file_name=' + file_name,
            type: "post",
            beforeSend: function () {
                results = []
                status_list = []
                $("#table_result > tbody").html('');
                $('#row_detail').hide();
                $("#row_results").hide();
            console.log("[ML SERVICE] process_doc")
            },
        }).done(function (jsondata, textStatus, jqXHR) {
            console.log(" JSON0 --- /process_doc = ", jsondata)
            task_id = jsondata['go_task_id']
            status = jsondata['status']
            console.log("go_task_id = ", task_id)

            status_list.push(task_id)

            // for (i = 0; i < jsondata.length; i++) {
            //     task_id = jsondata['go_task_id']
            //     status = jsondata['status']
            //     console.log("task_id = ", task_id)

            //     // results.push(URL_GO + jsondata[i]['url_result'])
            //     status_list.push(task_id)
            //     result_button = `<button class="btn btn-small btn-success" style="display: none" id="btn-view-get-images" data=${i}>View</a>`
            //     $("#table_result > tbody").append(`<tr><td>${task_id}</td><td id=${task_id}>${status}</td><td>${result_button}</td></tr>`);
            //     $("#row_results").show();
            // }

            // Phase 2: Probe requested Go Service status

            var interval = setInterval(refresh, 1000);
            var fake_progress = 0
            function refresh() {
                n_success = 0
                console.log("status_list.lenght", status_list.length)
                for (i = 0; i < status_list.length; i++) {
                    $.ajax({
                        url: URL_GO_KEYWORD_STATUS + status_list[i],
                        success: function (data) {
                            console.log("SUCCEED query ",URL_GO_KEYWORD_STATUS + status_list[i])
                            console.log("/process_doc data = ",data)
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
        // Phase 1: Request Go Service
        var prompt = $('#keyword_prompt').val();
        var query = encodeURIComponent(prompt);

        $.ajax({
            url: URL_GO + '/get_images?query=' + query,
            type: "post",
            beforeSend: function () {
                results = []
                status_list = []
                $("#table_result > tbody").html('');
                $('#row_detail').hide();
                $("#row_results").hide();
            console.log("Requested ML Service")
            },
        }).done(function (jsondata, textStatus, jqXHR) {
            console.log("Resp (/get_images) = ", jsondata)
            task_id = jsondata['go_task_id']
            status = jsondata['status']
            status_list.push(task_id)
            // for (i = 0; i < jsondata.length; i++) {
            //     task_id = jsondata['go_task_id']
            //     status = jsondata['status']
            //     // results.push(URL_GO + jsondata[i]['url_result'])
            //     status_list.push(task_id)
            //     result_button = `<button class="btn btn-small btn-success" style="display: none" id="btn-view-get-images" data=${i}>View</a>`
            //     $("#table_result > tbody").append(`<tr><td>${task_id}</td><td id=${task_id}>${status}</td><td>${result_button}</td></tr>`);
            //     $("#row_results").show();
            // }

            // Phase 2: Probe requested Go Service status

            var interval = setInterval(refresh, 1000);
            var fake_progress = 0
            function refresh() {
                n_success = 0
                for (i = 0; i < status_list.length; i++) {
                    $.ajax({
                        url:  URL_GO_IMAGE_STATUS + status_list[i],
                        success: function (data) {
                            console.log("SUCCEED query get_image Status on GO")
                            console.log("data = ",data)
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

    $('#btn-get-health').on('click', function () {
       // Phase 1: Request Go Service
       $.ajax({
           url: URL_GO + '/health?msg=ping',
           type: "post",
           beforeSend: function () {
               results = []
               status_list = []
               $("#table_result > tbody").html('');
               $('#row_detail').hide();
               $("#row_results").hide();
           console.log("Checking Go Health")
           },
           success: function (jsondata, textStatus, jqXHR) {
               console.log("GO Response = ", jsondata);
               console.log('jqXHR = ', jqXHR);
           },
           error: function (jqXHR, textStatus, errorThrown) {
               console.log("Failed Requesting service at Go, probaly Go Dead")
               console.log(jqXHR);
               console.log("error = ", errorThrown)
           },
       })
    })

    $('#btn-test').on('click', function () {
        // Phase 1: Request Go Service
        $.ajax({
            url: URL_GO + '/test?query=hi&task_id=',
            type: "post",
            beforeSend: function () {
                results = []
                status_list = []
                $("#table_result > tbody").html('');
                $('#row_detail').hide();
                $("#row_results").hide();
            console.log("Checking Go Health")
            },
            success: function (jsondata, textStatus, jqXHR) {
                console.log("GO Response = ", jsondata);
                console.log('jqXHR = ', jqXHR);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.log("Failed Requesting service at Go, probaly Go Dead")
                console.log(jqXHR);
                console.log("error = ", errorThrown)
            },
        }).done(function (jsondata, textStatus, jqXHR) {
            console.log("[SUCCEED] Requested Test endpont GO!");
            console.log("jqXHR status",jqXHR.status);
            // jsondata.length == number of download pictures (tasks)
            for (i = 0; i < jsondata.length; i++) {
                task_id = jsondata[i]['task_id']
                status = jsondata[i]['status']
                results.push(URL_VLNML + jsondata[i]['url_result'])
                status_list.push(task_id)
                result_button = `<button class="btn btn-small btn-success" style="display: none" id="btn-view-detect" data=${i}>View</a>`
                $("#table_result > tbody").append(`<tr><td>${task_id}</td><td id=${task_id}>${status}</td><td>${result_button}</td></tr>`);
                $("#row_results").show();
            }

            // Phase 2: Probe requested ML Service status
            var interval = setInterval(refresh, 1000);
            var fake_progress = 0;

            function refresh() {
                n_success = 0
                for (i = 0; i < status_list.length; i++) {
                    // Loop through each task
                    $.ajax({
                        url: URL_VLNML_STATUS + status_list[i],
                        success: function (data,textStatus,jqXHR) {
                            id = status_list[i]
                            status = data['status']
                            console.log("[+] Task ",i," -- ",id)
                            console.log("[+] status = ", status)
                            console.log("[+] data = ", data)
                            console.log("[+] jqXHR = ", jqXHR)
                            $('#' + id).html(status)
                            if ((status == 'SUCCESS') || (status == 'FAILED')) {
                                // show View Button for each Completed Task
                                console.log("[+] show View Button")
                                console.log("[+] jqXHR = ", jqXHR)
                                console.log("[+]  $($('#' + id) = ", $($('#' + id)))
                                console.log("[+]  $($('#' + id).siblings()[1]) = ", $($('#' + id)).siblings()[1])
                                console.log("[+]  $($('#' + id).siblings()[1]).children() = ", $($('#' + id).siblings()[1]).children())
                                $($('#' + id).siblings()[1]).children().show()
                                n_success++
                            }

                            // Show progress bar
                            console.log("fake_progress",fake_progress)
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
                    // Complete prob status for ALL tasks
                    $('#progress-bar').hide();
                    clearInterval(interval);
                }
            }
        }).fail(function (jsondata, textStatus, jqXHR) {
            console.log(jsondata)
            $("#row_results").hide();
        });
     })
    // View Buttons
    $(document).on('click', '#btn-view-detect', function (e) {
        id = $(e.target).attr('data')
        console.log("View-Detect Button clicked")

        console.log("e = ", e)
        console.log("id = ", id)
        console.log("results = ", results)
        $.get(results[id], function (data) {
            console.log("data = ", data)
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
                url: URL_VLNML_STATUS + status_list[i],
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
