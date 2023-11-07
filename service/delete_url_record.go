package service

func (s *UrlService) deleteUrlRecord(shortUrl string) error {
	return s.UrlRecordDao.Delete(shortUrl)
}
