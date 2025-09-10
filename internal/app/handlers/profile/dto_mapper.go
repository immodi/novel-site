package profile

import profiledtostructs "immodi/novel-site/internal/http/structs/profile"

func MapToProfile(username, pictureUrl, joinDate string) *profiledtostructs.ProfileDto {
	return &profiledtostructs.ProfileDto{
		Username:   username,
		PictureURL: pictureUrl,
		JoinDate:   joinDate,
	}
}
