import React from 'react';
import {connect} from "dva";
import styles from './index.module.css';
import { Descriptions, Spin } from "antd";
import {VERSION} from "@/utils/About"

@connect(({ loading, user }) => {
    return {
      settingAbout: user.settingAbout,
      settingAboutLoading: loading.effects['user/getSettingAbout'],
    }
 })

class About extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'user/getSettingAbout' 
    });
  }

  render() {
    const {settingAbout, settingAboutLoading} = this.props;
    return (
      <Spin spinning={settingAboutLoading}>
        <div className={styles.title}>关于</div>
        <Descriptions column={1}>
          <Descriptions.Item label="操作系统">{settingAbout['SystemInfo']}</Descriptions.Item>
          <Descriptions.Item label="Golang版本">{settingAbout['Golangversion']}</Descriptions.Item>
          <Descriptions.Item label="Gin版本">{settingAbout['GinVersion']}</Descriptions.Item>
          <Descriptions.Item label="Spug 版本">{VERSION}</Descriptions.Item>
          <Descriptions.Item label="官网文档">
            <a href="https://spug.dev" target="_blank" rel="noopener noreferrer">https://spug.dev</a>
          </Descriptions.Item>
        </Descriptions>
      </Spin>
    )
  }
}

export default About
