package cmd

/*
	func  TestPush(c *C) {
		url, clean := s.TemporalDir()
		defer clean()

		server, err := PlainInit(url, true)
		c.Assert(err, IsNil)

		_, err = s.Repository.CreateRemote(&config.RemoteConfig{
			Name: "test",
			URLs: []string{url},
		})
		c.Assert(err, IsNil)

		err = s.Repository.Push(&PushOptions{
			RemoteName: "test",
		})
		c.Assert(err, IsNil)

		AssertReferences(c, server, map[string]string{
			"refs/heads/master": "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			"refs/heads/branch": "e8d3ffab552895c19b9fcf7aa264d277cde33881",
		})

		AssertReferences(c, s.Repository, map[string]string{
			"refs/remotes/test/master": "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			"refs/remotes/test/branch": "e8d3ffab552895c19b9fcf7aa264d277cde33881",
		})
	}
func TemporalHomeDir() (path string, clean func()) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	fs := osfs.New(home)
	relPath, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}

	path = fs.Join(fs.Root(), relPath)
	clean = func() {
		_ = billyutil.RemoveAll(fs, relPath)
	}

	return
}
*/
