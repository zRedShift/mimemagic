package mimemagic

var treeMagicSignatures = []treeMagic{
	{997, []treeMatch{{"VIDEO_TS/VIDEO_TS.IFO", -1, fileType, false, false, false, nil}, {"VIDEO_TS/VIDEO_TS.IFO;1", -1, fileType, false, false, false, nil}, {"VIDEO_TS.IFO", -1, fileType, false, false, false, nil}, {"VIDEO_TS.IFO;1", -1, fileType, false, false, false, nil}}},
	{998, []treeMatch{{"HVDVD_TS/HV000I01.IFO", -1, fileType, false, false, false, nil}, {"HVDVD_TS/HV001I01.IFO", -1, fileType, false, false, false, nil}, {"HVDVD_TS/HVA00001.VTI", -1, fileType, false, false, false, nil}}},
	{995, []treeMatch{{".autorun", -1, fileType, true, false, false, nil}, {"autorun", -1, fileType, true, false, false, nil}, {"autorun.sh", -1, fileType, true, false, false, nil}}},
	{996, []treeMatch{{"BDAV", -1, directoryType, false, false, true, nil}, {"BDMV", -1, directoryType, false, false, true, nil}}},
	{985, []treeMatch{{"AUDIO_TS/AUDIO_TS.IFO", -1, fileType, false, false, false, nil}, {"AUDIO_TS/AUDIO_TS.IFO;1", -1, fileType, false, false, false, nil}}},
	{991, []treeMatch{{".kobo", -1, directoryType, false, false, true, nil}, {"system/com.amazon.ebook.booklet.reader", -1, anyType, false, false, false, nil}}},
	{1001, []treeMatch{{"autorun.exe", -1, fileType, false, true, false, nil}, {"autorun.inf", -1, fileType, false, false, false, nil}}},
	{992, []treeMatch{{"dcim", -1, directoryType, false, false, true, nil}}},
	{993, []treeMatch{{"PICTURES", -1, directoryType, true, false, true, nil}}},
	{999, []treeMatch{{"MPEG2/AVSEQ01.MPG", -1, fileType, false, false, false, nil}}},
	{1000, []treeMatch{{"mpegav/AVSEQ01.DAT", -1, fileType, false, false, false, nil}}},
}
