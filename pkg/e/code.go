package e

const (
	// 与HTTP返回码对应
	CodeCheckLogin      = 401
	CodeNoPermissionErr = 403
	CodeNotFound        = 404

	// 通用错误
	CodeDBError    = 10001
	CodeParamError = 10002

	// 用户相关
	CodeLanguageSet = 20001
	CodeValidateError = 20002

	// 文件相关
	CodeCreateFolderFailed   = 30001
	CodeCheckChunk           = 30002
	CodeErrGetUploadChunk    = 30003
	CodeUploadChunk          = 30004
	CodeErrMerge             = 30005
	CodeErrRename            = 30006
	CodeErrDelete            = 30007
	CodeErrSetShare          = 30008
	CodeErrGetMyFiles        = 30009
	CodeErrDownload          = 30010
	CodeDownloadFileNotExist = 30011
	// admin
	CodeListUserErr  = 40001
	CodeSetStatusErr = 40002
	CodeGetFilesErr  = 40003
	CodeReviewErr    = 40004
	// 未定错误
	CodeNotSet = -1
)
