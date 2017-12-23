(function () {
    $.ajax({
        type: 'POST',
        url: contextPath + "/machines",
        success: function (content, textStatus, request) {
            var machinesHtml = createMachinesTable(content)
            $('#machines').html(machinesHtml)
        },
        error: function (jqXHR, textStatus, errorThrown) {
            alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
        }
    })


    function createMachinesTable(content) {
        var machinesHtml = '<table>' +
            '<tr>' +
            '<td>Name</td>' +
            '<td>State</td>' +
            '<td>Hostname</td>' +
            '<td>IP</td>' +
            '<td>OS</td>' +
            '<td>Cores</td>' +
            '<td>Total Memory</td>' +
            '<td>Available Memory</td>' +
            '<td>Memory Used</td>' +
            '<td>Total Disk</td>' +
            '<td>Free Disk</td>' +
            '<td>Disk Used</td>' +
            '<td>Cost</td>' +
            '</tr>'
        if (content && content.length) {
            for (var j = 0; j < content.length; j++) {
                var machineResult = content[j]
                var hardwareInfo = machineResult.MachineInfo

                machinesHtml += '<tr><td class="MachineName">' + machineResult.Name + '</td>'
                    + '<td>' + (machineResult.Error || 'OK') + '</td>'
                    + '<td>' + hardwareInfo.Hostname + '</td>'
                    + '<td>' + (hardwareInfo.Ips && hardwareInfo.Ips.length > 0 ? hardwareInfo.Ips.join(', ') : '') + '</td>'
                    + '<td>' + hardwareInfo.OS + '</td>'
                    + '<td>' + hardwareInfo.Cores + '</td>'
                    + '<td>' + hardwareInfo.HumanizedTotalMemory + '</td>'
                    + '<td>' + hardwareInfo.HumanizedAvailableMemory + '</td>'
                    + '<td>' + new Number(hardwareInfo.MemoryUsedPercent).toFixed(2) + '%</td>'
                    + '<td>' + hardwareInfo.HumanizedTotalDisk + '</td>'
                    + '<td>' + hardwareInfo.HumanizedFreeDisk + '</td>'
                    + '<td>' + new Number(hardwareInfo.DiskUsedPercent).toFixed(2) + '%</td>'
                    + '<td>' + machineResult.CostTime + '</td>'
                    + '</tr>'
            }
        }

        machinesHtml += '</table>'
        return machinesHtml;
    }

})()