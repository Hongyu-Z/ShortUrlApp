package service

func (s *UrlService) deleteUrlRecord(shortUrl string) error {
	return s.urlRecordDao.Delete(shortUrl)
}
