(function () {
    $.ajax({
        type: 'POST',
        url: contextPath + "/machines",
        success: function (content, textStatus, request) {
            var machines = $('#machines')
            var machinesHtml = '<table>' +
                '<tr><td>Name</td><td>State</td><td>Hostname</td><td>OS</td><td>Total Memory</td><td>Free Memory</td>' +
                '<td>Memory Used</td><td>Cores</td><td>Total Disk</td><td>Free Disk</td><td>Disk Used</td></tr>'
            if (content && content.length) {
                for (var j = 0; j < content.length; j++) {
                    var machineResult = content[j]
                    var hardwareInfo = machineResult.MachineInfo
                    machinesHtml += '<tr><td>' + machineResult.Name + '</td>'
                        + '<td>' + (machineResult.Error || 'OK') + '</td>'
                        + '<td>' + hardwareInfo.Hostname + '</td>'
                        + '<td>' + hardwareInfo.OS + '</td>'
                        + '<td>' + hardwareInfo.HumanizedTotalMemory + '</td>'
                        + '<td>' + hardwareInfo.HumanizedFreeMemory + '</td>'
                        + '<td>' + new Number(hardwareInfo.MemoryUsedPercent).toFixed(2) + '%</td>'
                        + '<td>' + hardwareInfo.Cores + '</td>'
                        + '<td>' + hardwareInfo.HumanizedTotalDisk + '</td>'
                        + '<td>' + hardwareInfo.HumanizedFreeDisk + '</td>'
                        + '<td>' + new Number(hardwareInfo.DiskUsedPercent).toFixed(2) + '%</td>'
                        + '</tr>'
                }
            }

            machinesHtml += '</table>'
            machines.html(machinesHtml)
        },
        error: function (jqXHR, textStatus, errorThrown) {
            alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
        }
    })


})()